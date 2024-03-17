package parse_form

import (
	"fmt"
	"os"
	"testing"
)

// TestParseTable2D 生成二维数组, 类似课程表样子的
func TestParseTable2D(t *testing.T) {
	table, err := os.ReadFile("Pg课程html/TestForm.html")
	if err != nil {
		t.Error(err)
	}
	parseTable, err := ParseTable2D(&table)
	if err != nil {
		t.Error(err)
	}
	TableStand := [12][9]string{
		{"上午", "一", "", "", "", "", "", "", ""},
		{"上午", "二", "", "", "", "", "", "", ""},
		{"上午", "三", "", "", "高级算法设计与分析1班｛12-19周[教师:孙冬璞,地点:待定]｝", "", "", "", ""},
		{"上午", "四", "", "", "高级算法设计与分析1班｛12-19周[教师:孙冬璞,地点:待定]｝", "", "", "", ""},
		{"下午", "五", "高级算法设计与分析1班｛12-19周[教师:孙冬璞,地点:待定]｝", "机器学习1班｛3-10周[教师:陈晨,地点:待定]｝", "", "", "", "", ""},
		{"下午", "六", "高级算法设计与分析1班｛12-19周[教师:孙冬璞,地点:待定]｝", "机器学习1班｛3-10周[教师:陈晨,地点:待定]｝", "", "", "", "", ""},
		{"下午", "七", "", "组合数学1班｛4-11周[教师:高峻,地点:待定]｝", "", "机器学习1班｛3-10周[教师:陈晨,地点:待定]｝", "", "", ""},
		{"下午", "八", "", "组合数学1班｛4-11周[教师:高峻,地点:待定]｝", "", "机器学习1班｛3-10周[教师:陈晨,地点:待定]｝", "", "", ""},
		{"晚上", "九", "", "", "", "组合数学1班｛4-11周[教师:高峻,地点:待定]｝", "", "", ""},
		{"晚上", "十", "", "", "", "组合数学1班｛4-11周[教师:高峻,地点:待定]｝", "", "", ""},
		{"晚上", "十一", "", "", "", "", "", "", ""},
		{"晚上", "十二", "", "", "", "", "", "", ""},
	}
	for i := 1; i < 13; i++ {
		for j := 1; j < 10; j++ {
			if parseTable[i][j] != TableStand[i-1][j-1] {
				t.Errorf("parseTable[%d][%d] = %s, want %s", i, j, parseTable[i][j], TableStand[i-1][j-1])
			}
		}
	}

	//fmt.Print("end\n\n")
	//for i := 1; i < 13; i++ {
	//	for j := 1; j < 10; j++ {
	//		if parseTable[i][j] == "" {
	//			fmt.Print("\" \",")
	//		} else {
	//			fmt.Print("\"" + parseTable[i][j] + "\", ")
	//		}
	//	}
	//	fmt.Println()
	//}
	//生成标准输出表格
}

// TestParseTable1D 传入某一周课表, 生成课程信息切片, 研究生
func TestParseTable1D(t *testing.T) {
	table, err := os.ReadFile("Pg课程html/TestForm.html")
	if err != nil {
		t.Error(err)
	}
	parseTable, err := ParseTable1D(&table, 2)
	if err != nil {
		t.Error(err)
	}
	for _, v := range parseTable {
		fmt.Println(v)
	}
}

// TestParseTablePgAll 给定一个本学期课表, 生成课程信息切片, 研究生
func TestParseTablePgAll(t *testing.T) {
	table, err := os.ReadFile("Pg课程html/TestForm.html")
	if err != nil {
		t.Error(err)
	}
	parseTable, err := ParseTablePgAll(&table)
	if err != nil {
		t.Error(err)
	}
	for _, v := range parseTable {
		fmt.Println(v)
	}
}

// TestParseTablePgByWeek 给定一个学期的课表, 生成课程信息切片, 解析本科生的
func TestParseTableUgAll(t *testing.T) {
	table, err := os.ReadFile("Ug课程html/本学期课表1.html")
	if err != nil {
		t.Error(err)
	}
	allCoures, err := ParseTableUgAll(&table)
	if err != nil {
		t.Error(err)
	}
	for _, course := range allCoures {
		fmt.Println(course)
	}
	return
}

// TestParseTableUgByWeek 给定一个学期的课表和某一周, 生成那周的课程信息切片, 解析本科生的
func TestParseTableUgByWeek(t *testing.T) {
	table, err := os.ReadFile("Ug课程html/本学期课表.html")
	if err != nil {
		t.Error(err)
	}
	allCoures, err := ParseTableUgByWeek(&table, 1)
	if err != nil {
		t.Error(err)
	}
	for _, course := range allCoures {
		fmt.Println(course)
	}
	return
}

// TestParse_json_ug_nd 解析农大课表json
func TestParse_json_ug_nd(t *testing.T) {
	jsonInfo, err := os.ReadFile("农大课表json/tsconfig.json")
	if err != nil {
		t.Error(err)
	}
	res, err := Parse_json_ug_nd(&jsonInfo)
	if err != nil {
		t.Error(err)
	}
	for _, v := range res {
		fmt.Println(v)
	}
}

// 解析哈理工本科成绩
func TestParseTableUgSoreHrbust(t *testing.T) {
	table, err := os.ReadFile("Ug成绩/Ug成绩.html")
	if err != nil {
		t.Error(err)
	}
	allCoures, err := ParseTableUgSore(&table, "2023", "1")
	if err != nil {
		t.Error(err)
	}
	for _, course := range allCoures {
		fmt.Println(course)
	}
	return
}
