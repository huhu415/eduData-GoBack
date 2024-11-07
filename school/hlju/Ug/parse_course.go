package hljuUg

import (
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"eduData/repository"
	"eduData/school/pub"

	"github.com/sirupsen/logrus"
)

/*
"KCWZSM": null,
"RWH": "2024-2025-1-1412016010-01",
"SFFXEXW": null,
"FILEURL": null,
"SKSJ": "计算机控制\n[于海华]\n[1-9周][3号楼-512]\n第3-4节",
"XB": 1,
"SKSJ_EN": "Computer Control \n[于海华]\n[1-9Week][3号楼-512]\n3-4",
"KEY": "xq2_jc2" or "bz
*/
type schedule struct {
	KCWZSM  *string `json:"KCWZSM"`  // 课程位置说明，可以为 null
	RWH     string  `json:"RWH"`     // 任务号，唯一标识课程安排
	SFFXEXW *string `json:"SFFXEXW"` // 是否为选修课，可以为 null
	FILEURL *string `json:"FILEURL"` // 课程相关文件的 URL，可以为 null
	SKSJ    string  `json:"SKSJ"`    // 课程时间和详情（中文）
	XB      int     `json:"XB"`      // 相同xb, 为同一门课程
	SKSJ_EN *string `json:"SKSJ_EN"` // 课程时间和详情（英文），可以为 null
	KEY     string  `json:"KEY"`     // 格式: 星期?_节次? or bz备注
}

func ParseCoruse(data *[]byte) ([]repository.Course, error) {
	var schedules []schedule
	if err := json.Unmarshal(*data, &schedules); err != nil {
		return nil, err
	}

	// xb排序
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].XB < schedules[j].XB
	})

	// 初始化颜色队列, 用于给课程上色, 一共19个, 应该用不完
	preXB := -999
	color := ""
	queue := pub.NewColorList()

	var resCoures []repository.Course
	for _, schedule := range schedules {
		information := schedule.SKSJ

		if schedule.KEY == "bz" {
			pattern := `(?P<CourseName>.+?) \[(?P<Weeks>\d+-\d+周)\] (?P<Teacher>\S+) 备注:(?P<Remark>.*)`
			re := regexp.MustCompile(pattern)
			if match := re.FindStringSubmatch(information); match != nil {
				// courseName, weeks, teacher, remark := match[1], match[2], match[3], match[4]
				courseName, weeks, teacher, _ := match[1], match[2], match[3], match[4]

				startWeek, endWeek, _, err := pub.ExtractWeekRange(weeks)
				if err != nil {
					logrus.Errorf("解析周次失败: %s", err)
					return nil, err
				}
				course := repository.Course{
					School:                "hlju",
					CourseContent:         courseName,
					TeacherName:           teacher,
					Week:                  0,
					WeekDay:               0,
					NumberOfLessons:       0,
					NumberOfLessonsLength: 0,
					BeginWeek:             startWeek,
					EndWeek:               endWeek,
				}
				resCoures = append(resCoures, course)
			} else {
				logrus.Warnf("课程信息解析失败: %s", information)
				return nil, errors.New("课程信息解析失败")
			}
		} else {
			result := strings.Split(information, "\n")
			logrus.Debugf("result: %v, result.len(): %d", result, len(result))
			switch len(result) {
			case 4:
				//[高级语言程序设计（C语言） [马慧] [5-15周][汇文楼-437] 第9-10节]
				className, teacherName, weeksAndlocation, Time := result[0], result[1], result[2], result[3]
				// 提取周次
				startWeek, endWeek, _, err := pub.ExtractWeekRange(weeksAndlocation)
				if err != nil {
					logrus.Warnf("提取周次失败, result-%v: %s", result, err)
					end := strings.Index(weeksAndlocation, "]")
					re := regexp.MustCompile(`\d+`)
					logrus.Debugf("weeksAndlocation[:end]: %s", weeksAndlocation[:end])
					matches := re.FindAllString(weeksAndlocation[:end], -1)

					switch len(matches) {
					case 1:
						num, err := strconv.Atoi(matches[0])
						if err != nil {
							return nil, err
						}
						startWeek, endWeek = num, num
					case 2:
						num, err := strconv.Atoi(matches[0])
						if err != nil {
							return nil, err
						}
						startWeek = num

						num, err = strconv.Atoi(matches[1])
						if err != nil {
							return nil, err
						}
						endWeek = num
					}

				}

				// 提取地点
				location := ""
				start := strings.LastIndex(weeksAndlocation, "[")
				end := strings.LastIndex(weeksAndlocation, "]")
				if start != -1 && end != -1 && start < end && strings.Count(weeksAndlocation, "[") >= 2 {
					location = weeksAndlocation[start+1 : end]
				} else {
					logrus.Warnf("解析地点失败: %s, 没有第二个[]", weeksAndlocation)
				}

				// 提取课程节数
				startCourse, endCourse, err := pub.ExtractCoruse(Time)
				if err != nil {
					logrus.Errorf("解析课程节数失败: %s", err)
					return nil, err
				}
				logrus.Debugf("time: %s", Time)
				logrus.Debugf("课程节数: %d-%d", startCourse, endCourse)

				// 提取第几大节课***
				// todo 观望一下, 没问题就可以删了
				var bigCourse int
				reBigCourse := regexp.MustCompile(`jc(\d+)`)
				matchbigCourse := reBigCourse.FindStringSubmatch(schedule.KEY)
				if len(matchbigCourse) > 1 {
					bigCourse, _ = strconv.Atoi(matchbigCourse[1])
					logrus.Debugf("jc 后面的数字是: %s, conv: %d\n", matchbigCourse[1], bigCourse*2-1)
					if bigCourse*2-1 != startCourse {
						logrus.Warnf("jc 后面的数字和课程节数不匹配")
					}
				}

				// 提取星期几***
				var weekDay int
				reWeekDay := regexp.MustCompile(`xq(\d+)`)
				matchweekDay := reWeekDay.FindStringSubmatch(schedule.KEY)
				if len(matchweekDay) > 1 {
					weekDay, err = strconv.Atoi(matchweekDay[1])
					if err != nil {
						logrus.Warn("没找到xq后面的数字")
					}
					logrus.Debugf("xq 后面的数字是: %s, conv: %d\n", matchweekDay[1], weekDay)
				}

				// 如果和上一次的XB不一样, 那么就从队列中取出新的颜色
				if schedule.XB != preXB {
					preXB = schedule.XB
					color = queue.Remove(queue.Front()).(string)
				}

				for i := startWeek; i <= endWeek; i++ {
					course := repository.Course{
						School:                "hlju",
						CourseContent:         className,
						TeacherName:           teacherName,
						CourseLocation:        location,
						NumberOfLessons:       startCourse,
						NumberOfLessonsLength: endCourse - startCourse + 1,
						BeginWeek:             startWeek,
						EndWeek:               endWeek,
						Week:                  i,
						WeekDay:               weekDay,
						Color:                 color,
					}
					resCoures = append(resCoures, course)
				}
			case 5:
				// [【2024秋期初补缓考考试】 运动控制 9月27日 18:00-20:00 汇文楼-325]
				name, content, date, time, location := result[0], result[1], result[2], result[3], result[4]

				// 提取星期几***
				var err error
				var weekDay int
				reWeekDay := regexp.MustCompile(`xq(\d+)`)
				matchweekDay := reWeekDay.FindStringSubmatch(schedule.KEY)
				if len(matchweekDay) > 1 {
					weekDay, err = strconv.Atoi(matchweekDay[1])
					if err != nil {
						logrus.Warn("没找到xq后面的数字")
					}
					logrus.Debugf("xq 后面的数字是: %s, conv: %d\n", matchweekDay[1], weekDay)
				}

				// 提取第几大节课***
				// todo 观望一下, 没问题就可以删了
				var bigCourse int
				reBigCourse := regexp.MustCompile(`jc(\d+)`)
				matchbigCourse := reBigCourse.FindStringSubmatch(schedule.KEY)
				if len(matchbigCourse) > 1 {
					bigCourse, err = strconv.Atoi(matchbigCourse[1])
					if err != nil {
						logrus.Warn("没找到jc后面的数字")
					}
					logrus.Debugf("jc 后面的数字是: %s, conv: %d\n", matchbigCourse[1], bigCourse*2-1)
				}

				course := repository.Course{
					School:                "hlju",
					CourseContent:         name + content,
					TeacherName:           date + time,
					CourseLocation:        location,
					WeekDay:               weekDay,
					NumberOfLessons:       bigCourse*2 - 1,
					NumberOfLessonsLength: bigCourse * 2,
					// Week                  int
					// Color                 string
				}
				resCoures = append(resCoures, course)
			default:
				logrus.Warnf("不存在的课程长度: %s", information)
				return nil, errors.New("不存在的课程长度")
			}
		}
	}

	return resCoures, nil
}
