package neauUg

import (
	"fmt"
	"os"
	"testing"
)

const (
	// 东北农业大学本科生账号密码
	USERNAME string = "A08220441"
	PASSWORD string = "shiyunxin0527@"
)

// TestSigninUgNEAU 测试东北农业大学本科生登陆
func TestUgSignin(t *testing.T) {
	_, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Error(err)
	}
}

/*----------------------------------------------------------------------*/

// TestRevLeftChidScorePg 东北农业大学本科生, 接收成绩json
func TestGetJSONneau(t *testing.T) {
	cookiejar, err := Signin(USERNAME, PASSWORD)
	if err != nil {
		t.Error(err)
	}
	jsoNneau, err := GetData(cookiejar, "2024-2025-1-1")
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*jsoNneau))
	}
}

/*----------------------------------------------------------------------*/

// TestParse_json_ug_nd 解析农大课表json
func TestParse_json_ug_nd(t *testing.T) {
	jsonInfo, err := os.ReadFile("农大课表json/tsconfig.json")
	if err != nil {
		t.Error(err)
	}
	res, err := ParseData(&jsonInfo)
	if err != nil {
		t.Error(err)
	}
	for _, v := range res {
		fmt.Println(v)
	}
}
