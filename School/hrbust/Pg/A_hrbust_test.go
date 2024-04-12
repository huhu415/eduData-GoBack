// 哈理工研究生相关测试
package hrbustPg

import (
	"fmt"
	"os"
	"testing"
)

const (
	// 研究生账号密码
	USERNAME string = "2320410125"
	PASSWORD string = "Aa788415"

	// 研究生不同的登陆页面
	//登陆后成功后跳转的页面
	DEFAULT string = "Default.aspx?UID="
	//左边菜单
	LEFTMENU string = "leftmenu.aspx?UID="
	//学期课表
	LEFTTERM string = "Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID="
	//某一周课表
	LEFETTHISWEEK string = "Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID="
)

// TestSignin 测试哈尔滨理工大学研究生登陆
func TestSignin(t *testing.T) {
	// 测试登陆
	_, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("登陆失败: %s", err)
		return
	}
}

/*----------------------------------------------------------------------*/

// default渲染后, 调出Leftmenu和topmenu, Leftmenu中有很多chid, 其中有课表, 成绩等.

// TestGetData 研究生发送LEFTTERM, 获得学期课表
func TestGetData(t *testing.T) {
	cookiejar, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	CourseTable, err := GetData(cookiejar, USERNAME, LEFTTERM)
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*CourseTable))
	}
}

// TestGetData2 研究生发送LEFETTHISWEEK, 获得某一周课表, 并且可以选择某一周
func TestGetData2(t *testing.T) {
	cookiejar, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	CourseTable, err := GetData(cookiejar, USERNAME, LEFETTHISWEEK, "DropDownListWeeks=DropDownListWeeks=13")
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*CourseTable))
	}
}

/*----------------------------------------------------------------------*/

// TestParseDataCoures1D 传入某一周课表, 生成课程信息切片, 研究生
func TestParseDataCoures1D(t *testing.T) {
	table, err := os.ReadFile("Pg课程html/TestForm.html")
	if err != nil {
		t.Error(err)
	}
	parseTable, err := ParseDataCoures1D(&table, 2)
	if err != nil {
		t.Error(err)
	}
	for _, v := range parseTable {
		fmt.Println(v)
	}
}

// TestParseDataCoures2D 生成二维数组, 类似课程表样子的
func TestParseDataCoures2D(t *testing.T) {
	table, err := os.ReadFile("Pg课程html/TestForm.html")
	if err != nil {
		t.Error(err)
	}
	parseTable, err := ParseDataCoures2D(&table)
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

// TestParseDataCouresAll 给定一个本学期课表, 生成课程信息切片, 研究生
func TestParseDataCouresAll(t *testing.T) {
	table, err := os.ReadFile("Pg课程html/TestForm.html")
	if err != nil {
		t.Error(err)
	}
	parseTable, err := ParseDataCouresAll(&table)
	if err != nil {
		t.Error(err)
	}
	for _, v := range parseTable {
		fmt.Println(v)
	}
}
