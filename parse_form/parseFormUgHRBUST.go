package parse_form

import (
	"container/list"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"eduData/database"
)

func chineseToNumber(chinese rune) (int, error) {
	switch chinese {
	case '一':
		return 1, nil
	case '二':
		return 2, nil
	case '三':
		return 3, nil
	case '四':
		return 4, nil
	case '五':
		return 5, nil
	case '六':
		return 6, nil
	case '日':
		return 7, nil
	default:
		return 0, fmt.Errorf("不支持的汉字数字: %s", string(chinese))
	}
}

// ExtractWeekRange 提取周数范围, 返回起始周, 结束周, 单双周, 单双周默认是5, 0是双周, 1是单周, 便于运算
func ExtractWeekRange(text string) (startWeek, endWeek, evenOrOdd int, err error) {
	startWeek, endWeek, evenOrOdd = 0, 0, 5
	if []rune(text)[0] == '第' {
		// 匹配形式 : 第3周
		// 正则表达式提取数字
		WeekSinge := regexp.MustCompile("[0-9]+").FindAllString(text, 1)
		var atoi int
		atoi, err = strconv.Atoi(WeekSinge[0])
		if err != nil {
			err = fmt.Errorf("形式 第x周 无法解析起始周: %v", err)
			return
		}
		//如果是这种形式, 那么起始周与结束周都是这个数字
		startWeek, endWeek = atoi, atoi
		//fmt.Println("起始周-结束周", startWeek, "-", endWeek)
	} else {
		// 匹配形式 : 1-15周 或 1-15单周
		matchWeekRange := regexp.MustCompile(`(\d+)-(\d+)`).FindStringSubmatch(text)
		if len(matchWeekRange) == 0 {
			err = errors.New("无法匹配周数范围, x-x")
			return
		}
		// 第一个数
		startWeek, err = strconv.Atoi(matchWeekRange[1])
		if err != nil {
			err = fmt.Errorf("形式 x-x周 无法解析起始周: %v", err)
			return
		}
		// 第二个数
		endWeek, err = strconv.Atoi(matchWeekRange[2])
		if err != nil {
			err = fmt.Errorf("形式 x-x周 无法解析结束周: %v", err)
			return
		}
		// 判断单双周
		switch {
		case strings.Contains(text, "单"):
			evenOrOdd = 1
		case strings.Contains(text, "双"):
			evenOrOdd = 0
		}
		//fmt.Println("起始周-结束周", startWeek, "-", endWeek)
	}
	return
}

// ParseTableUgAll 给定一个学期的课表, 返回这个学期的所有课程, 解析本科生的
func ParseTableUgAll(table *[]byte) ([]database.Course, error) {
	//创建返回的变量
	var courses []database.Course

	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(*table)))
	if err != nil {
		return nil, err
	}

	//判断是否能找到课程信息
	if doc.Find("table.infolist_tab tbody tr.infolist_common").Length() == 0 {
		return nil, errors.New("not find table.infolist_hr tbody tr.infolist_common")
	}

	//创建一个课程
	course := database.Course{
		// 因为函数是解析本科生, 所以这是本科生
		StuType: 1,
		School:  "hrbust",
	}
	// 初始化颜色队列, 用于给课程上色, 一共19个, 应该用不完
	queue := list.New()
	queue.PushBack("#2e1f54")
	queue.PushBack("#52057f")
	queue.PushBack("#bf033b")
	queue.PushBack("#f00a36")
	queue.PushBack("#ff6908")
	queue.PushBack("#ffc719")
	queue.PushBack("#598c14")
	queue.PushBack("#335238")
	queue.PushBack("#4a8594")
	queue.PushBack("#051736")
	queue.PushBack("#706357")
	queue.PushBack("#b0a696")
	queue.PushBack("#004eaf")
	queue.PushBack("#444444")
	queue.PushBack("#c1d1e0")
	queue.PushBack("#c1d1e0")
	queue.PushBack("#faa918")
	queue.PushBack("#8f1010")
	queue.PushBack("#d2ea32")

	// 遍历每个课程
	doc.Find("table.infolist_tab tbody tr.infolist_common").Each(func(trIndex int, row *goquery.Selection) {
		// 遍历每个课程的元素, 课程名称, 任教老师, 上课时间、地点等
		row.Find("td").Each(func(tdIndex int, cell *goquery.Selection) {
			// 课程名称
			if tdIndex == 2 {
				course.CourseContent = strings.TrimSpace(cell.Text())
			}
			// 任课教师
			if tdIndex == 3 {
				course.TeacherName = strings.TrimSpace(cell.Text())
			}
			// 上课时间、地点			或单纯的汉字
			if tdIndex == 9 {
				if cell.Find("table.none tbody tr").Length() != 0 {
					// 如果有课程时间的table, 把颜色掏出来
					color := queue.Remove(queue.Front()).(string)

					// 开始遍历个课程中上课时间地点的每行
					cell.Find("table.none tbody tr").Each(func(trIndexIn int, cellIn *goquery.Selection) {
						startWeek, endWeek := 0, 0
						// 单周或双周, 或者没有
						evenOrOdd := 0

						// 开始遍历一行的每个元素
						cellIn.Find("td").Each(func(tdIndexIn int, cellInIn *goquery.Selection) {
							text := strings.TrimSpace(cellInIn.Text())

							// 第0个是周数
							if tdIndexIn == 0 {
								startWeek, endWeek, evenOrOdd, err = ExtractWeekRange(text)
								if err != nil {
									return
								}
							}

							// 1是星期几
							if tdIndexIn == 1 {
								if text == "" {
									course.WeekDay = 0
								} else {
									course.WeekDay, err = chineseToNumber([]rune(text)[2])
									if err != nil {
										fmt.Println(err)
										return
									}
								}
							}

							// 2是第几节
							if tdIndexIn == 2 {
								switch text {
								case "一上午":
									course.NumberOfLessons = 1
									course.NumberOfLessonsLength = 4
								case "一下午":
									course.NumberOfLessons = 5
									course.NumberOfLessonsLength = 4
								case "第一大节":
									course.NumberOfLessons = 1
									course.NumberOfLessonsLength = 2
								case "第二大节":
									course.NumberOfLessons = 3
									course.NumberOfLessonsLength = 2
								case "第三大节":
									course.NumberOfLessons = 5
									course.NumberOfLessonsLength = 2
								case "第四大节":
									course.NumberOfLessons = 7
									course.NumberOfLessonsLength = 2
								case "第五大节":
									course.NumberOfLessons = 9
									course.NumberOfLessonsLength = 2
								case "第六大节":
									course.NumberOfLessons = 11
									course.NumberOfLessonsLength = 2
								default:
									course.NumberOfLessons = 0
									course.NumberOfLessonsLength = 0
								}
							}

							// 3是地点
							if tdIndexIn == 3 {
								// 把全角括号(中文)替换成半角括号(英文)
								replaced := strings.ReplaceAll(text, "（", "(")
								replaced = strings.ReplaceAll(replaced, "）", ")")
								course.CourseLocation = replaced
							}
						})

						// 一个课程的其中一行结束, 准备添加到切片结构体里
						if course.NumberOfLessons == 0 || course.NumberOfLessonsLength == 0 || course.WeekDay == 0 {
							course.Week = 0
							course.BeginWeek = startWeek
							course.EndWeek = endWeek
							courses = append(courses, course)
							// 没有时间或地点的课程, 就不用上色了, 有默认颜色#c1d1e0
						} else {
							//根据单双周, 添加到切片结构体里, 并且有课程, 需要上色
							course.Color = color
							for i := startWeek; i <= endWeek; i++ {
								// 如果是单双周, 则判断是否符合
								if evenOrOdd != 5 {
									if i%2 != evenOrOdd {
										continue
									}
								}
								course.Week = i
								course.BeginWeek = startWeek
								course.EndWeek = endWeek
								courses = append(courses, course)
							}
						}
					})
				} else {
					// 如果没有课程时间的table, 那么就是单纯的汉字, 或空白
					text := strings.TrimSpace(cell.Text())
					// 判断空白与否
					if len(text) == 0 {
						// 空白的就不放数据库里了
						course.WeekDay, course.Week, course.NumberOfLessons, course.NumberOfLessonsLength = 0, 0, 0, 0
						course.CourseLocation = ""
					} else {
						// 类似形式 : 1-15周 时间地点都不占
						startWeek, endWeek, _, err := ExtractWeekRange(text)
						if err != nil {
							return
						}
						course.BeginWeek, course.EndWeek = startWeek, endWeek
						course.WeekDay, course.Week, course.NumberOfLessons, course.NumberOfLessonsLength = 0, 0, 0, 0
						courses = append(courses, course)
					}
				}
			}
		})
	})

	// 如果解析有问题, 返回错误
	if err != nil {
		return nil, err
	}
	return courses, nil
}

// ParseTableUgByWeek 给定一个学期的课表和某一周, 返回这个学期的这周的课程, 解析本科生的
func ParseTableUgByWeek(table *[]byte, week int) ([]database.Course, error) {
	//创建返回的变量
	var courses []database.Course

	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(*table)))
	if err != nil {
		return nil, err
	}

	//判断是否能找到课程信息
	if doc.Find("table.infolist_tab tbody tr.infolist_common").Length() == 0 {
		return nil, errors.New("not find table.infolist_hr tbody tr.infolist_common")
	}

	doc.Find("table.infolist_tab tbody tr.infolist_common").Each(func(trIndex int, row *goquery.Selection) {
		course := database.Course{
			Week: week,
		}
		row.Find("td").Each(func(tdIndex int, cell *goquery.Selection) {
			if tdIndex == 2 {
				course.CourseContent = strings.TrimSpace(cell.Text())
			}
			if tdIndex == 3 {
				course.TeacherName = strings.TrimSpace(cell.Text())
			}
			if tdIndex == 9 {
				// 课程时间
				cell.Find("table.none tbody tr").Each(func(trIndexIn int, cellIn *goquery.Selection) {
					isWeekIn := false
					cellIn.Find("td").Each(func(tdIndexIn int, cellInIn *goquery.Selection) {
						// 第0个是周数
						if tdIndexIn == 0 {
							text := strings.TrimSpace(cellInIn.Text())
							if []rune(text)[0] == '第' {
								// 匹配形式 : 第3周
								// 正则表达式提取数字
								WeekSinge := regexp.MustCompile("[0-9]+").FindAllString(text, 1)
								atoi, err := strconv.Atoi(WeekSinge[0])
								if err != nil {
									fmt.Println(err)
									return
								}
								fmt.Printf("第%d周\n", atoi)
								if atoi == week {
									isWeekIn = true
								}
							} else {
								// 匹配形式 : 1-15周
								// 删除‘周’
								text = strings.ReplaceAll(text, "周", "")
								// 分割字符串, 然后取第一个和第二个数字
								weekRangeParts := strings.Split(text, "-")
								startWeek, err := strconv.Atoi(weekRangeParts[0])
								if err != nil {
									fmt.Println("无法解析起始周:", err)
									return
								}
								endWeek, err := strconv.Atoi(weekRangeParts[1]) // 去掉周字
								if err != nil {
									fmt.Println("无法解析结束周:", err)
									return
								}
								fmt.Println("起始周-结束周", startWeek, "-", endWeek)
								if startWeek <= week && week <= endWeek {
									isWeekIn = true
								}
							}
						}
						if isWeekIn {
							text := strings.TrimSpace(cellInIn.Text())
							// 1是星期几
							if tdIndexIn == 1 {
								course.WeekDay, err = chineseToNumber([]rune(text)[2])
								if err != nil {
									fmt.Println(err)
									return
								}
							}
							// 2是第几节
							if tdIndexIn == 2 {
								course.NumberOfLessons, err = chineseToNumber([]rune(text)[1])
								if err != nil {
									fmt.Println(err)
									return
								}
								if []rune(text)[2] == '大' {
									course.NumberOfLessons = 2
								} else {
									course.NumberOfLessons = 1
								}
							}
							// 3是地点
							if tdIndexIn == 3 {
								course.CourseLocation = strings.TrimSpace(cellInIn.Text())
							}
						}
					})
					if isWeekIn {
						courses = append(courses, course)
					}
				})
			}
		})
		fmt.Println()
	})
	return courses, nil
}
