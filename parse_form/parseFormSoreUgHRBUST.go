package parse_form

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"eduData/database"
)

func ParseTableUgSore(table *[]byte) ([]database.Course, error) {
	// 使用 goquery 解析 HTML 表格
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(*table)))
	if err != nil {
		return nil, err
	}

	//判断是否能找到课程信息
	if doc.Find("table.datalist tbody tr").Length() == 0 {
		return nil, errors.New("not find table.datalist tbody tr")
	}
	doc.Find("table.datalist tbody tr").Each(func(trIndex int, row *goquery.Selection) {
		row.Find("td").Each(func(tdIndex int, cell *goquery.Selection) {
			// 课程名称
			if tdIndex == 3 {
				// 课程名称
				// course.Name = cell.Text()
			}
			if tdIndex == 6 {
				// 总评
				// course.Name = cell.Text()
			}
			if tdIndex == 7 {
				//学分
				// course.Name = cell.Text()
			}
			if tdIndex == 8 {
				//学时
				// course.Name = cell.Text()
			}
			if tdIndex == 9 {
				// 任课属性
				// course.Name = cell.Text()
			}
		})
	})
	return nil, nil
}
