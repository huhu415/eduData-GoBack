// Package hrbust 哈尔滨理工大学
package hrbustPg

import (
	"fmt"
	"io"
	"net/http/cookiejar"

	"net/http"

	"eduData/setting"
)

const (
	HTTPHOST string = "http://yjs.hrbust.edu.cn/Gstudent/"
	HOST     string = "yjs.hrbust.edu.cn"
)

// GetData 用于接收左边的各个子成员, 子成员在RevLeftmenu()的返回值里面找
// 例如LeftChidUrl = Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID=
// CookieAppend = DropDownListWeeks=DropDownListWeeks=13*/
func GetData(cookie *cookiejar.Jar, username, LeftChidUrl string, CookieAppend ...any) (*[]byte, error) {
	// 从setting中获取UserAgent
	var userAgent = setting.UserAgent

	//解析参数
	var CookieAppendRes string
	for _, v := range CookieAppend {
		CookieAppendRes += "; " + v.(string)
	}

	//新建一个客户端
	client := &http.Client{
		Jar: cookie,
	}

	req, err := http.NewRequest("GET", HTTPHOST+LeftChidUrl+username, nil)
	if err != nil {
		return nil, err
	}
	// 注释取消就可以获取某一周课表, 但cookie要是string的
	//req.Header.Add("cookie", cookie+"; LoginType=LoginType=1"+CookieAppendRes)
	req.Header.Add("Host", HOST)
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", HTTPHOST)
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	//发送请求, 并接收响应, 同时defer关闭响应体

	readAllBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(readAllBody))

	return &readAllBody, nil
}
