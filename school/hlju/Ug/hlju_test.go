package hljuUg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	USERNAME = ""
	PASSWORD = ""
)

func TestSignin(t *testing.T) {
	assert := assert.New(t)

	cookie, err := Signin(USERNAME, PASSWORD)
	if assert.Nil(err, "登陆失败") {
		t.Run("GetData-ThisTerm", func(t *testing.T) {
			// 意思为2024-2025年 1学期
			CourseTable, err := GetData(cookie, 2024, 1)
			if assert.Nil(err, "获取课表失败") {
				t.Log(string(*CourseTable))
			}
		})
	}
}

/*----------------------------------------------------------------------*/

func TestParseDataCoures(t *testing.T) {
	assert := assert.New(t)

	table, err := os.ReadFile("html/ug.html")
	if assert.Nil(err, "读取文件失败") {
		parseTable, err := ParseCoruse(&table)
		if assert.Nil(err, "解析课表失败") {
			for _, v := range parseTable {
				t.Log(v)
			}
		}
	}
}
