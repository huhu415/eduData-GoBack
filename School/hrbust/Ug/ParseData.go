package hrbustUg

import (
	"eduData/School/pub"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"eduData/models"
)

// ParseDataCrouseAll 给定一个学期的课表, 返回这个学期的所有课程, 解析本科生的
func ParseDataCrouseAll(table *[]byte) ([]models.Course, error) {
	//创建返回的变量
	var courses []models.Course

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
	course := models.Course{
		// 因为函数是解析本科生, 所以这是本科生
		StuType: 1,
		School:  "hrbust",
	}
	// 初始化颜色队列, 用于给课程上色, 一共19个, 应该用不完
	queue := pub.NewColorList()

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
								startWeek, endWeek, evenOrOdd, err = pub.ExtractWeekRange(text)
								if err != nil {
									return
								}
							}

							// 1是星期几
							if tdIndexIn == 1 {
								if text == "" {
									course.WeekDay = 0
								} else {
									course.WeekDay, err = pub.ChineseToNumber([]rune(text)[2])
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
					if len(text) == 0 {
						// 空白的
						course.BeginWeek, course.EndWeek = 0, 0
						course.CourseLocation = ""
					} else {
						// 类似形式 : 1-15周 时间地点都不占
						startWeek, endWeek, _, err := pub.ExtractWeekRange(text)
						if err != nil {
							return
						}
						course.BeginWeek, course.EndWeek = startWeek, endWeek
					}
					course.WeekDay, course.Week, course.NumberOfLessons, course.NumberOfLessonsLength = 0, 0, 0, 0
					courses = append(courses, course)
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

// ParseDataCrouseByWeek 给定一个学期的课表和某一周, 返回这个学期的这周的课程, 解析本科生的
func ParseDataCrouseByWeek(table *[]byte, week int) ([]models.Course, error) {
	//创建返回的变量
	var courses []models.Course

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
	course := models.Course{
		// 因为函数是解析本科生, 所以这是本科生
		StuType: 1,
		School:  "hrbust",
		Week:    week,
	}

	// 初始化颜色队列, 用于给课程上色, 一共19个, 应该用不完
	queue := pub.NewColorList()

	doc.Find("table.infolist_tab tbody tr.infolist_common").Each(func(trIndex int, row *goquery.Selection) {
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
								startWeek, endWeek, evenOrOdd, err = pub.ExtractWeekRange(text)
								if err != nil {
									return
								}
							}

							// 1是星期几
							if tdIndexIn == 1 {
								if text == "" {
									course.WeekDay = 0
								} else {
									course.WeekDay, err = pub.ChineseToNumber([]rune(text)[2])
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
							if week >= startWeek && week <= endWeek {
								if evenOrOdd != 5 {
									// 如果是单双周, 则判断是否符合
									if week%2 == evenOrOdd {
										course.BeginWeek = startWeek
										course.EndWeek = endWeek
										courses = append(courses, course)
									}
								} else {
									course.BeginWeek = startWeek
									course.EndWeek = endWeek
									courses = append(courses, course)
								}
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
						startWeek, endWeek, _, err := pub.ExtractWeekRange(text)
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
		// debug用
		//fmt.Println()
	})
	return courses, nil
}

// ParseDataSore 解析哈理工本科生成绩页面
func ParseDataSore(table *[]byte, year, term string) ([]models.CourseGrades, error) {
	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(*table)))
	if err != nil {
		return nil, err
	}

	var courseGrades []models.CourseGrades
	//判断是否能找到课程信息
	//if doc.Find("table.datalist tbody tr").Length() == 0 {
	//	return nil, errors.New("not find table.datalist tbody tr")
	//}

	courseGrade := models.CourseGrades{
		// 哈理工本科生
		School:   "hrbust",
		StuType:  1,
		Year:     year,
		Semester: term,
	}

	// 每个课程的循环
	doc.Find("table.datalist tbody tr").Each(func(trIndex int, row *goquery.Selection) {
		// 每个属性的循环
		row.Find("td").Each(func(tdIndex int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			// 课程名称
			if tdIndex == 3 {
				courseGrade.CourseName = text
			}

			// 总评
			if tdIndex == 6 {
				grade, err := strconv.ParseFloat(text, 64)
				if err != nil {
					switch text {
					case "优秀":
						grade = 95
					case "良好":
						grade = 85
					case "中":
						grade = 75
					case "及格":
						grade = 65
					case "不及格":
						grade = 0
					}
				}
				courseGrade.CourseGrade = grade
			}

			//学分
			if tdIndex == 7 {
				credit, err := strconv.ParseFloat(text, 64)
				if err != nil {
					credit = 0
				}
				courseGrade.CourseCredit = credit
			}

			////学时
			//if tdIndex == 8 {
			//
			//}

			// 任课属性
			if tdIndex == 9 {
				courseGrade.CourseType = text

				// 一行有消息全部获取完毕, 加入切片后结束
				courseGrades = append(courseGrades, courseGrade)
				return
			}
		})
	})
	return courseGrades, nil
}
