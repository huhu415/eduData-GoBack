// 哈理工本科生相关测试
package hrbustUg

import (
	"eduData/bootstrap"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// 本科生账号密码
	USERNAME string = "2204010417"
	PASSWORD        = "13737826060a"
	YEARTEST        = "43"
	TERMTEST        = "2"
)

func TestHrbustUg(t *testing.T) {
	bootstrap.Loadconfig()
	assert := assert.New(t)

	cookie, err := Signin(USERNAME, PASSWORD)
	if assert.Nil(err, "登陆失败") {
		// 获取学期课表
		t.Run("GetData-ThisTerm", func(t *testing.T) {
			CourseTable, err := GetData(cookie, "2000")
			if assert.Nil(err, "获取课表失败") {

				// 解析课表
				t.Run("ParseDataCrouseAll", func(t *testing.T) {
					allCoures, err := ParseDataCrouseAll(CourseTable)
					if assert.Nil(err, "解析课程失败") {
						for _, course := range allCoures {
							t.Log(course)
						}
					}
				})

			}
		})

		// 获取成绩
		t.Run("GetDataScore", func(t *testing.T) {
			Score, err := GetDataScore(cookie, YEARTEST, TERMTEST)
			if assert.Nil(err, "获取成绩失败") {

				// 解析成绩
				t.Run("ParseDataSore", func(t *testing.T) {
					allCoures, err := ParseDataSore(Score, YEARTEST, TERMTEST)
					if assert.Nil(err, "解析成绩失败") {
						for _, course := range allCoures {
							t.Log(course)
						}
					}
				})

			}
		})

	}
}

/*----------------------------------------------------------------------*/

// TestParseTablePgByWeek 给定一个学期的课表, 生成课程信息切片, 解析本科生的
func TestParseTableUgAll(t *testing.T) {
	// table, err := os.ReadFile("Ug课程html/本学期课表.html")
	table, err := os.ReadFile("Ug课程html/currcourse.html")
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
