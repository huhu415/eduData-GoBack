package parse_form

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"eduData/database"
)

// ParseTableUgSore 解析本科生成绩页面
func ParseTableUgSore(table *[]byte, year, term string) ([]database.CourseGrades, error) {
	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(*table)))
	if err != nil {
		return nil, err
	}

	var courseGrades []database.CourseGrades
	//判断是否能找到课程信息
	//if doc.Find("table.datalist tbody tr").Length() == 0 {
	//	return nil, errors.New("not find table.datalist tbody tr")
	//}

	courseGrade := database.CourseGrades{
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
					case "中等":
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
