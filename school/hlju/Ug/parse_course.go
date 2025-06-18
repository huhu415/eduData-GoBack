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

	"github.com/huhu415/gorange"
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
			// 官方备注信息, 意思是显示在官方网页中
			pattern := `(?P<CourseName>.+?) \[(?P<Weeks>.*周)\] (?P<Teacher>\S+) 备注:(?P<Remark>.*)`
			re := regexp.MustCompile(pattern)
			match := re.FindStringSubmatch(information)
			if match == nil {
				logrus.Errorf("备注信息解析失败: %s", information)
				continue
			}

			// courseName, weeks, teacher, remark := match[1], match[2], match[3], match[4]
			courseName, weeks, teacher, _ := match[1], match[2], match[3], match[4]
			wRange, err := gorange.ExtractRange(weeks)
			if err != nil {
				logrus.Errorf("解析课程节数失败: %s", err)
				return nil, err
			}
			startWeek, endWeek := wRange[0], wRange[len(wRange)-1]

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
			continue
		}

		result := strings.Split(information, "\n")
		logrus.Debugln()
		logrus.Debugf("result: %v, result.len(): %d", result, len(result))

		if !strings.Contains(information, "][") {
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
				logrus.Debugf("考试: jc 后面的数字是: %s, conv: %d\n", matchbigCourse[1], bigCourse*2-1)
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
			continue
		}

		logrus.Debug(result)
		// [大学体育A（Ⅰ） 身体素质练习 [韩洋] [5-17周][田径场-田径场（2号田径场）] 第5-6节]  len-5
		// [高级语言程序设计（C语言） [马慧] [5-15周][汇文楼-437] 第9-10节]					len-4
		// [新生研讨课【01】 [关秀娟] [8周][汇文楼-466]]									len-3
		className, teacherName, weeksAndlocation, Time := "", "", "", ""
		lenResult := len(result)
		switch {
		case lenResult == 3:
			className, teacherName, weeksAndlocation = result[lenResult-3], result[lenResult-2], result[lenResult-1]
			Time = "" // 如果只有3个元素，Time设为空字符串
		case lenResult >= 4:
			className, teacherName, weeksAndlocation, Time = result[lenResult-4], result[lenResult-3], result[lenResult-2], result[lenResult-1]
		default:
			return nil, errors.New("课程信息解析失败")
		}
		t := ""
		for i := 0; i < len(result)-4; i++ {
			t += result[i]
		}
		className = t + className
		logrus.Debugf("课程名称:%s, 老师名称%s", className, teacherName)

		// 提取周次, [5-15周][汇文楼-437]中的第一个[]中的内容
		end := strings.Index(weeksAndlocation, "]")
		weeks := weeksAndlocation[:end+1]
		weekRange, err := gorange.ExtractRange(weeks)
		if err != nil {
			logrus.Warnf("提取num-num周次失败, weeksAndlocation is %v, err:%s", weeks, err)
			continue
		}
		logrus.Debugf("最终提取出来的rangeWeek, %v", weekRange)

		// 提取地点
		location := ""
		start := strings.LastIndex(weeksAndlocation, "[")
		end = strings.LastIndex(weeksAndlocation, "]")
		if start != -1 && end != -1 && start < end && strings.Count(weeksAndlocation, "[") >= 2 {
			location = weeksAndlocation[start+1 : end]
			logrus.Debugf("提取出来的地点: %s", location)
		} else {
			logrus.Warnf("解析地点失败: %s, 没有第二个[]", weeksAndlocation)
		}

		// 提取课程节数
		startCourse, endCourse := 0, 0
		if Time == "" {
			// // 提取第几大节课***
			// todo 观望一下, 没问题就可以删了
			var bigCourse int
			reBigCourse := regexp.MustCompile(`jc(\d+)`)
			matchbigCourse := reBigCourse.FindStringSubmatch(schedule.KEY)
			if len(matchbigCourse) > 1 {
				bigCourse, _ = strconv.Atoi(matchbigCourse[1])
				logrus.Debugf("jc 后面的数字是: %s, conv: %d\n", matchbigCourse[1], bigCourse*2-1)
				startCourse = bigCourse*2 - 1
				endCourse = bigCourse * 2
			}
		} else {
			courseRange, err := gorange.ExtractRange(Time)
			if err != nil {
				logrus.Errorf("解析课程节数失败: %s", err)
				return nil, err
			}
			startCourse, endCourse = courseRange[0], courseRange[len(courseRange)-1]
			logrus.Debugf("待解析节数: %s, 解析后: %d-%d", Time, startCourse, endCourse)
		}

		// 提取星期几***
		var weekDay int
		reWeekDay := regexp.MustCompile(`xq(\d+)`)
		matchweekDay := reWeekDay.FindStringSubmatch(schedule.KEY)
		if len(matchweekDay) > 1 {
			weekDay, _ = strconv.Atoi(matchweekDay[1])
			logrus.Debugf("schedule.KEY: %s, 提取出来的是星期:%d\n", schedule.KEY, weekDay)
		} else {
			logrus.Warn("没找到xq后面的数字, 也就是说没有星期几")
		}

		if weekDay == 0 || (endCourse == 0 && startCourse == 0) || startCourse == 0 {
			beginWeek, endWeek := arryMaxMin(weekRange)

			course := repository.Course{
				School:                "hlju",
				CourseContent:         className,
				TeacherName:           teacherName,
				CourseLocation:        location,
				NumberOfLessons:       startCourse,
				NumberOfLessonsLength: endCourse - startCourse + 1,
				BeginWeek:             beginWeek,
				EndWeek:               endWeek,
				Week:                  0,
				WeekDay:               weekDay,
			}
			resCoures = append(resCoures, course)

			continue
		}

		// 如果和上一次的XB不一样, 那么就从队列中取出新的颜色
		if schedule.XB != preXB {
			preXB = schedule.XB
			color = queue.Remove(queue.Front()).(string)
		}

		for _, i := range weekRange {
			course := repository.Course{
				School:                "hlju",
				CourseContent:         className,
				TeacherName:           teacherName,
				CourseLocation:        location,
				NumberOfLessons:       startCourse,
				NumberOfLessonsLength: endCourse - startCourse + 1,
				// BeginWeek:             startWeek,
				// EndWeek:               endWeek,
				Week:    i,
				WeekDay: weekDay,
				Color:   color,
			}
			resCoures = append(resCoures, course)
		}
	}

	for i := range resCoures {
		resCoures[i].CourseContent = strings.TrimSpace(resCoures[i].CourseContent)
		resCoures[i].CourseContent = pub.FullWidthToHalfWidth(resCoures[i].CourseContent)
	}
	return resCoures, nil
}

func arryMaxMin(arr []int) (int, int) {
	if len(arr) == 0 {
		return 0, 0
	}

	// 创建副本以免修改原数组
	sorted := make([]int, len(arr))
	copy(sorted, arr)

	sort.Ints(sorted)

	return sorted[0], sorted[len(sorted)-1]
}
