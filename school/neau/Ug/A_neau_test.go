package neauUg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// 东北农业大学本科生账号密码
	USERNAME     = "A08220441"
	PASSWORD     = "shiyunxin0527@"
	NEAUDATATEST = "2024-2025-1-1"
)

func TestNeauUg(t *testing.T) {
	assert := assert.New(t)

	cookiejar, err := Signin(USERNAME, PASSWORD)
	if assert.Nil(err, "登陆失败") {

		// 获取课表
		t.Run("GetData", func(t *testing.T) {
			jsoNneau, err := GetData(cookiejar, NEAUDATATEST)
			assert.Nilf(err, "获取json失败")

			t.Log(string(*jsoNneau))
		})
	}

}

/*----------------------------------------------------------------------*/

// TestParse_json_ug_nd 解析农大课表json
func TestParse_json_ug_nd(t *testing.T) {
	jsonInfo, err := os.ReadFile("农大课表json/tsconfig.json")
	assert.Nilf(t, err, "读取json文件失败")

	res, err := ParseData(&jsonInfo)
	assert.Nilf(t, err, "解析json失败")

	for _, v := range res {
		t.Log(v)
	}
}
