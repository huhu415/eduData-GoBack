// 哈理工本科生相关测试
package hrbustUg

import (
	"fmt"
	"os"
	"testing"
)

const (
	// 本科生账号密码
	USERNAME string = "2204010417"
	PASSWORD string = "13737826060a"
)

// TestSingInUg 测试哈尔滨理工大学本科生登陆
func TestSignin(t *testing.T) {
	_, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("登陆失败: %s", err)
	}
}

/*----------------------------------------------------------------------*/

// TestRevLeftChidUg 本科生, 接收点击左侧地址后的html
func TestGetData(t *testing.T) {
	// TestRevLeftChid 发送LEFTTERM, 获得学期课表
	cookiejar, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	CourseTable, err := GetData(cookiejar, "2000")
	if err != nil {
		t.Errorf("error: %s", err)
		//t.Fatalf("error: %s", err)
	} else {
		fmt.Println(string(*CourseTable))
	}
}

// TestRevLeftChidScoreUg 本科生, 接收成绩html
func TestGetDataScore(t *testing.T) {
	cookiejar, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	Score, err := GetDataScore(cookiejar, "43", "2")
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*Score))
	}
}

/*----------------------------------------------------------------------*/

// TestParseTablePgByWeek 给定一个学期的课表, 生成课程信息切片, 解析本科生的
func TestParseTableUgAll(t *testing.T) {
	table, err := os.ReadFile("Ug课程html/本学期课表.html")
	if err != nil {
		t.Error(err)
	}
	allCoures, err := ParseDataCrouseAll(&table)
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
	allCoures, err := ParseDataCrouseByWeek(&table, 1)
	if err != nil {
		t.Error(err)
	}
	for _, course := range allCoures {
		fmt.Println(course)
	}
	return
}

// 解析哈理工本科成绩
func TestParseDataSore(t *testing.T) {
	table, err := os.ReadFile("Ug成绩/Ug成绩.html")
	if err != nil {
		t.Error(err)
	}
	allCoures, err := ParseDataSore(&table, "2023", "1")
	if err != nil {
		t.Error(err)
	}
	for _, course := range allCoures {
		fmt.Println(course)
	}
	return
}
