package hljuUg

import (
	"bytes"
	"eduData/models"
	"eduData/school/pub"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

type class struct {
	yearTerm  string  // 学年学期
	number    string  // 课程号
	name      string  // 课程名
	totalTime int     // 总学时
	credit    int     // 学分
	ctime     []ctime // 上课时间
	teacher   string  // 教师
	school    string  // 校区
}
type pair struct {
	first  int
	second int
}
type ctime struct {
	place        string // 上课地点
	weekDay      int    // 星期几
	sectionRange pair   // 第几节课
	weekRange    []pair // 周数范围
}

func ParseData(data *[]byte) ([]models.Course, error) {
	var courses []models.Course
	queue := pub.NewColorList()
	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*data))
	if err != nil {
		return nil, err
	}

	//判断是否能找到课程信息
	if doc.Find("table.ui_table").Length() == 0 {
		return nil, errors.New("not find table.ui_table")
	}

	doc.Find("table.ui_table").First().Find("tr.t_con").Each(func(i int, s *goquery.Selection) {
		class := class{}
		stackWeek := []pair{}
		s.Find("td").Each(func(j int, sm *goquery.Selection) {
			switch j {
			case 0:
				class.yearTerm = sm.Text()
			case 1:
				class.number = sm.Text()
			case 2:
				class.name = sm.Text()
			case 3:
				class.totalTime, _ = strconv.Atoi(sm.Text())
			case 4:
				class.credit, _ = strconv.Atoi(sm.Text())
			case 5:
				// 时间: 第几周, 星期几, 第几节课
				rawHTML, err := sm.Find("span").Html()
				if err != nil {
					logrus.Error(err)
				}
				rawHTML = strings.ReplaceAll(rawHTML, "\n", "")
				rawHTML = strings.ReplaceAll(rawHTML, "\t", "")
				rawHTML = strings.ReplaceAll(rawHTML, " ", "")
				lines := strings.Split(rawHTML, "<br/>")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line != "" && strings.Contains(line, "周") {
						matchWeekRange := regexp.MustCompile(`(\d+)-(\d+)周`).FindStringSubmatch(line)
						var startWeek, endWeek int
						startWeek, err = strconv.Atoi(matchWeekRange[1])
						if err != nil {
							fmt.Println("无法解析起始周:", err)
							return
						}
						endWeek, err = strconv.Atoi(matchWeekRange[2]) // 去掉周字
						if err != nil {
							fmt.Println("无法解析结束周:", err)
							return
						}
						// 遇到1-2周, 就放入栈中
						stackWeek = append(stackWeek, pair{startWeek, endWeek})
					}

					if matchWeek := regexp.MustCompile(`星期[一二三四五六日天]`).FindString(line); matchWeek != "" {
						week, err := pub.ChineseToNumber([]rune(matchWeek)[2])
						if err != nil {
							fmt.Println("无法解析星期:", err)
							return
						}

						// 解析从第几节到第几节
						matchSection := regexp.MustCompile(`(\d+)-(\d+)节`).FindStringSubmatch(line)
						var startSection, endSection int
						// 如果有节次
						if len(matchSection) >= 3 {
							startSection, err = strconv.Atoi(matchSection[1])
							if err != nil {
								fmt.Println("无法解析起始节:", err)
								return
							}
							endSection, err = strconv.Atoi(matchSection[2])
							if err != nil {
								fmt.Println("无法解析结束节:", err)
								return
							}
						} else {
							startSection = 0
							endSection = 0
						}

						// 将栈中的周数范围放入课程中
						class.ctime = append(class.ctime, ctime{
							weekRange:    stackWeek,
							weekDay:      week,
							sectionRange: pair{startSection, endSection},
						})
						// 如何有星期几了,清空栈
						stackWeek = make([]pair, 0)
					}
				}
			case 6:
				//地点
				rawHTML, err := sm.Find("span").Html()
				if err != nil {
					logrus.Error(err)
				}
				rawHTML = strings.ReplaceAll(rawHTML, "\n", "")
				rawHTML = strings.ReplaceAll(rawHTML, "\t", "")
				rawHTML = strings.ReplaceAll(rawHTML, " ", "")
				lines := strings.Split(rawHTML, "<br/>")
				for i, line := range lines {
					l := strings.TrimSpace(line)
					if l != "" {
						class.ctime[i].place = l
					}
				}
			case 7:
				//老师
				class.teacher = sm.Text()
			case 8:
				// 校本部
				class.school = sm.Text()
			}
		})

		color := queue.Remove(queue.Front()).(string)
		for _, c := range class.ctime {
			for _, weekRange := range c.weekRange {
				for week := weekRange.first; week <= weekRange.second; week++ {
					courses = append(courses, models.Course{
						StuID:                 "123123123123",
						School:                "hlju",
						Week:                  week,
						StuType:               1,
						WeekDay:               c.weekDay,
						NumberOfLessons:       c.sectionRange.first,
						NumberOfLessonsLength: c.sectionRange.second - c.sectionRange.first + 1,
						CourseContent:         class.name,
						Color:                 color,
						CourseLocation:        c.place,
						TeacherName:           class.teacher,
					})
				}
			}
		}

	})
	// for _, c := range courses {
	// 	fmt.Printf("%+v\n\n", c)
	// }
	return courses, nil
}
