package htmlgetter

import (
	"fmt"
	"testing"

	signin "eduData/sign_in"
)

const (
	// 本科生账号密码
	USERNAMEUG string = "2204010417"
	PASSWORDUG string = "13737826060a"
)

// TestRevLeftChidUg 本科生, 接收点击左侧地址后的html
func TestRevLeftChidUg(t *testing.T) {
	// TestRevLeftChid 发送LEFTTERM, 获得学期课表
	cookiejar, err := signin.SingInUg(USERNAMEUG, PASSWORDUG)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	CourseTable, err := RevLeftChidUg(cookiejar, "2000")
	if err != nil {
		t.Errorf("error: %s", err)
		//t.Fatalf("error: %s", err)
	} else {
		fmt.Println(string(*CourseTable))
	}
}

// TestRevLeftChidScoreUg 本科生, 接收成绩html
func TestRevLeftChidScoreUg(t *testing.T) {
	cookiejar, err := signin.SingInUg(USERNAMEUG, PASSWORDUG)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	Score, err := RevLeftChidScoreUg(cookiejar, "43", "2")
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*Score))
	}
}

const (
	// 研究生账号密码
	USERNAMEPG string = "2320410125"
	PASSWORDPG string = "Aa788415"

	// 研究生不同的登陆页面
	//登陆后成功后跳转的页面
	DEFAULT string = "Default.aspx?UID="
	//左边菜单
	LEFTMENU string = "leftmenu.aspx?UID="
	//学期课表
	LEFTTERM string = "Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID="
	//某一周课表
	LEFETTHISWEEK string = "Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID="

	// 本科生在函数中写死了

)

// default渲染后, 调出Leftmenu和topmenu, Leftmenu中有很多chid, 其中有课表, 成绩等.

// TestRevLeftChid 研究生发送LEFTTERM, 获得学期课表
func TestRevLeftChid(t *testing.T) {
	cookiejar, err := signin.SingInPg(USERNAMEPG, PASSWORDPG)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	CourseTable, err := RevLeftChidPg(cookiejar, USERNAMEPG, LEFTTERM)
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*CourseTable))
	}
}

// TestRevLeftChid2 研究生发送LEFETTHISWEEK, 获得某一周课表, 并且可以选择某一周
func TestRevLeftChid2(t *testing.T) {
	cookiejar, err := signin.SingInPg(USERNAMEPG, PASSWORDPG)
	if err != nil {
		t.Errorf("error: %s", err)
	}
	CourseTable, err := RevLeftChidPg(cookiejar, USERNAMEPG, LEFETTHISWEEK, "DropDownListWeeks=DropDownListWeeks=13")
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*CourseTable))
	}
}

const (
	// 东北农业大学本科生账号密码
	USERNAMEUGNEAU string = "A08220441"
	PASSWORDUGNEAU string = "shiyunxin0527@"
)

// TestRevLeftChidScorePg 东北农业大学本科生, 接收成绩json
func TestGetJSONneau(t *testing.T) {
	cookiejar, err := signin.SigninUgNEAU(USERNAMEUGNEAU, PASSWORDUGNEAU)
	if err != nil {
		t.Error(err)
	}
	jsoNneau, err := GetJSONneau(cookiejar, "2024-2025-1-1")
	if err != nil {
		t.Errorf("error: %s", err)
	} else {
		fmt.Println(string(*jsoNneau))
	}
}
