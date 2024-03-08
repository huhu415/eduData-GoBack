package sign_in

import (
	"testing"
)

// TestSingInUg 测试哈尔滨理工大学本科生登陆
func TestSingInUg(t *testing.T) {
	const (
		// 哈尔滨理工大学本科生 学号密码
		USERNAME string = "2204010417"
		PASSWORD string = "13737826060a"
	)
	_, err := SingInUg(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("登陆失败: %s", err)
	}
}

// TestSingInPg 测试哈尔滨理工大学研究生登陆
func TestSingInPg(t *testing.T) {
	const (
		// 哈尔滨理工大学研究生 学号密码
		USERNAME string = "2320410125"
		PASSWORD string = "Aa788415"
	)
	// 测试登陆
	_, err := SingInPg(USERNAME, PASSWORD)
	if err != nil {
		t.Errorf("登陆失败: %s", err)
		return
	}
}

// TestSigninUgNEAU 测试东北农业大学本科生登陆
func TestSigninUgNEAU(t *testing.T) {
	const (
		// 东北农业大学本科生 学号密码
		USERNAME string = "A08220441"
		PASSWORD string = "shiyunxin0527@"
	)
	_, err := SigninUgNEAU(USERNAME, PASSWORD)
	if err != nil {
		t.Error(err)
	}
}
