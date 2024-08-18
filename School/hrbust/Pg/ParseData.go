package hrbustPg

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/PuerkitoBio/goquery"

	"eduData/models"
)

// ParseDataCoures1D 要和Coures结构体配合使用, 传入某一周课表, 生成课程信息切片, 解析研究生的
// 我们前端的格式是星期几, 第几节, 上课长度, 上课内容, 上课地点, 可选参数第几周, int*/
func ParseDataCoures1D(table *[]byte, args ...any) ([]models.Course, error) {
	var courses []models.Course
	week := 0

	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*table))
	if err != nil {
		return nil, err
	}

	//判断是否能找到课程信息
	if doc.Find("table#contentParent_dgData tbody tr").Length() == 0 {
		return nil, errors.New("not find table#contentParent_dgData tbody tr")
	}

	//找到当前是第几周, 如果没有会取args[0]的值
	doc.Find("#contentParent_drpWeek_drpWeeks option").Each(func(rowIndex int, row *goquery.Selection) {
		selectedWeek, ok := row.Attr("selected")
		if ok && selectedWeek == "selected" {
			week, err = strconv.Atoi(row.Text())
			if err != nil {
				err = errors.New("weekDay is not int")
				return
			}
			return
		}
	})
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		week = args[0].(int)
	}
	if week == 0 {
		return nil, errors.New("not find the selected week in the table, you need add the weekDay by hand")
	}

	//课表不算表头, 有12行, 9列, 为什么有11列, 而不是9+1=10列, 是因为colIndexOwn要找到nil才能停下来, 如果是10列, 最后一个就是满的, 需要找到11列才能停下来
	var MapMaker [13][11]bool
	doc.Find("table#contentParent_dgData tbody tr").Each(func(rowIndex int, row *goquery.Selection) {
		colIndexOwn := 0 //如果换了一行, colIndexOwn要归零, 正常要从0开始, 但下面的函数已经完成colIndexOwn++了
		row.Find("td").Each(func(colIndex int, cell *goquery.Selection) {
			colIndexOwn++                                                        //colIndex在这里+1, 那么我们自己的的colIndexOwn也在这里+1
			Value := regexp.MustCompile(`\s+`).ReplaceAllString(cell.Text(), "") // 定义正则表达式，去除任何空白字符

			//fmt.Printf("Row: %d, Col: %d, rowspan: %s, Value: %s\n", rowIndex, colIndex, cell.AttrOr("rowspan", "0"), Value)
			//上面语句debug用

			//找到第一个没有被占用的位置, 来更新colIndexOwn, 因为colIndex有的时候会少, 因为合并的原因, 所以要一个标准的colIndexOwn
			for MapMaker[rowIndex][colIndexOwn] == true {
				colIndexOwn++
			}
			if len(Value) > 0 {
				//找不到rowspan, 返回值是0, 所以这里是能找到rowspan的话
				if v := cell.AttrOr("rowspan", "0"); v != "0" {
					rowspanNum, err := strconv.Atoi(v)
					if err != nil {
						fmt.Println(err)
					}

					//竖向做标记, 因为rowspanNum是跨行数.
					for i := rowspanNum; i > 0; i-- {
						MapMaker[rowIndex+i-1][colIndexOwn] = true
					}

					//要大于2个汉字的长度才能算是课程, 因为"晚上"不算课程
					if len(Value) > 8 {
						matchTeacher := regexp.MustCompile(`教师:([^\s,]+)`).FindAllStringSubmatch(Value, 1)
						matchLocation := regexp.MustCompile(`地点:([^\s\]]+)`).FindAllStringSubmatch(Value, 1)
						matchContent := regexp.MustCompile(`｛.*?｝`).ReplaceAllString(Value, "")
						courses = append(courses, models.Course{
							StuType:               2,
							Week:                  week,
							WeekDay:               colIndexOwn - 2, //因为前面还有2个td所以要减去2
							NumberOfLessons:       rowIndex,
							NumberOfLessonsLength: rowspanNum,
							CourseContent:         matchContent,
							TeacherName:           matchTeacher[0][1],
							CourseLocation:        matchLocation[0][1],
							//TeacherName:           "",
							//beginWeek:             0,
							//endWeek:               0,
						})
					}
				} else {
					MapMaker[rowIndex][colIndexOwn] = true
				}
			}
		})
	})
	//for _, row := range MapMaker {
	//	for _, col := range row {
	//		fmt.Print(col, " ")
	//	}
	//	fmt.Println()
	//}
	return courses, nil
}

// ParseDataCoures2D 传入某一周的带有table的html, 解析出[13][11]string的切片,
// 这里只针对课表, 其他表格不一定适用, 因为课表只有rowspan, 没有colspan
// 并且这个函数的结果是横坐标是星期1-8, 竖向坐标是1-12节课的内容string类型, 解析研究生的*/
func ParseDataCoures2D(table *[]byte) ([][]string, error) {
	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*table))
	if err != nil {
		return nil, err
	}
	if doc.Find("table#contentParent_dgData tbody tr").Length() == 0 {
		return nil, errors.New("not find table#contentParent_dgData tbody tr")
	}

	//课表不算表头, 有12行, 9列, 为什么有11列, 而不是9+1=10列, 是因为colIndexOwn要找到nil才能停下来, 如果是10列, 最后一个就是满的, 需要找到11列才能停下来
	MapMaker := make([][]string, 13)
	for i := 0; i < 13; i++ {
		MapMaker[i] = make([]string, 11)
	}
	//第一节课就是第一行,从左向右依次是时间是1, 节次是2,星期1是3, 到星期7是9. MapMaker从[1][1]开始.
	doc.Find("table#contentParent_dgData tbody tr").Each(func(rowIndex int, row *goquery.Selection) {
		colIndexOwn := 0 //如果换了一行, colIndexOwn要归零, 正常要从0开始, 但下面的函数已经完成colIndexOwn++了
		row.Find("td").Each(func(colIndex int, cell *goquery.Selection) {
			colIndexOwn++                                                        //colIndex在这里+1, 那么我们自己的的colIndexOwn也在这里+1
			Value := regexp.MustCompile(`\s+`).ReplaceAllString(cell.Text(), "") // 定义正则表达式，匹配任何空白字符

			//fmt.Printf("Row: %d, Col: %d, rowspan: %s, Value: %s\n", rowIndex, colIndex, cell.AttrOr("rowspan", "0"), Value)
			//上面语句debug用

			for MapMaker[rowIndex][colIndexOwn] != "" {
				colIndexOwn++
			}
			if v := cell.AttrOr("rowspan", "0"); v != "0" {
				//找不到rowspan, 返回值是0, 所以这里是能找到
				i, err := strconv.Atoi(v)
				if err != nil {
					fmt.Println(err)
				}
				for ; i > 0; i-- {
					MapMaker[rowIndex+i-1][colIndexOwn] = Value
				}
			} else {
				//如果能找不到rowspan, 直接赋值
				MapMaker[rowIndex][colIndexOwn] = Value
			}
		})
	})
	return MapMaker, err
}

// ParseDataCouresAll 传入html一个学期的课程表, 解析出带有所有课程的切片, 解析研究生的, 与其他的区别是, 这个会解析一学期的, 而不是某一周的
func ParseDataCouresAll(table *[]byte) ([]models.Course, error) {
	var courses []models.Course

	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*table))
	if err != nil {
		return nil, err
	}

	//判断是否能找到课程信息
	if doc.Find("table#contentParent_dgData tbody tr").Length() == 0 {
		return nil, errors.New("not find table#contentParent_dgData tbody tr")
	}

	//课表不算表头, 有12行, 9列, 为什么有11列, 而不是9+1=10列, 是因为colIndexOwn要找到nil才能停下来, 如果是10列, 最后一个就是满的, 需要找到11列才能停下来
	var MapMaker [13][11]bool
	doc.Find("table#contentParent_dgData tbody tr").Each(func(rowIndex int, row *goquery.Selection) {
		colIndexOwn := 0 //如果换了一行, colIndexOwn要归零, 正常要从0开始, 但下面的函数已经完成colIndexOwn++了
		row.Find("td").Each(func(colIndex int, cell *goquery.Selection) {
			colIndexOwn++                                                        //colIndex在这里+1, 那么我们自己的的colIndexOwn也在这里+1
			Value := regexp.MustCompile(`\s+`).ReplaceAllString(cell.Text(), "") // 定义正则表达式，去除任何空白字符

			//fmt.Printf("Row: %d, Col: %d, rowspan: %s, Value: %s\n", rowIndex, colIndex, cell.AttrOr("rowspan", "0"), Value)
			//上面语句debug用

			//找到第一个没有被占用的位置, 来更新colIndexOwn, 因为colIndex有的时候会少, 因为合并的原因, 所以要一个标准的colIndexOwn
			for MapMaker[rowIndex][colIndexOwn] == true {
				colIndexOwn++
			}
			if len(Value) > 0 {
				//找不到rowspan, 返回值是0, 所以这里是能找到rowspan的话
				if v := cell.AttrOr("rowspan", "0"); v != "0" {
					rowspanNum, _ := strconv.Atoi(v)

					//竖向做标记, 因为rowspanNum是跨行数.
					for i := rowspanNum; i > 0; i-- {
						MapMaker[rowIndex+i-1][colIndexOwn] = true
					}

					//要大于2个汉字的长度才能算是课程, 因为"晚上"不算课程
					if len(Value) > 8 {
						matchTeacher := regexp.MustCompile(`教师:([^\s,]+)`).FindAllStringSubmatch(Value, 1)
						matchLocation := regexp.MustCompile(`地点:([^\s\]]+)`).FindAllStringSubmatch(Value, 1)
						matchContent := regexp.MustCompile(`｛.*?｝`).ReplaceAllString(Value, "")
						course := models.Course{
							//id:                    1,
							//StuID:                 1,
							//Week:                  1,
							School:                "hrbust",
							StuType:               2,               //研究生
							WeekDay:               colIndexOwn - 2, //因为前面还有2个td所以要减去2
							NumberOfLessons:       rowIndex,
							NumberOfLessonsLength: rowspanNum,
							CourseContent:         matchContent,
							CourseLocation:        matchLocation[0][1],
							TeacherName:           matchTeacher[0][1],
						}
						matchWeekRange := regexp.MustCompile(`(\d+)-(\d+)周`).FindStringSubmatch(Value)
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
						for i := startWeek; i <= endWeek; i++ {
							course.Week = i
							courses = append(courses, course)
						}
					}
				} else {
					MapMaker[rowIndex][colIndexOwn] = true
				}
			}
		})
	})

	// 判断一下错误有没有传递出来
	if err != nil {
		return nil, err
	}
	return courses, nil
}
