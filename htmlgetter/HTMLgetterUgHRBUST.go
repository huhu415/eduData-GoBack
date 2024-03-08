// 本科生获取原始html课程

package htmlgetter

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"eduData/setting"
)

// RevLeftChidUg 获取原始的各种侧边栏内容, moduleId=2020是成绩查询, 2000是本学期课表
func RevLeftChidUg(cookieJar *cookiejar.Jar, moduleId string) (*[]byte, error) {
	// 从setting中获取UserAgent
	var userAgent = setting.UserAgent

	//新建一个客户端
	client := &http.Client{
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	//http://jwzx.hrbust.edu.cn/academic/accessModule.do?moduleId=2000&groupId=&randomString=20240221153427x91KU5
	//http://jwzx.hrbust.edu.cn/academic/accessModule.do?moduleId=2020&groupId=&randomString=20240304095909VGZjoo
	// 第三次请求, 请求课表
	newQuery := url.Values{}
	newQuery.Set("moduleId", moduleId)
	newQuery.Set("randomString", time.Now().Format("20060102150405")+strconv.Itoa(rand.Int()))
	req, err := http.NewRequest("GET", "http://jwzx.hrbust.edu.cn/academic/accessModule.do?"+newQuery.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//发送请求, 并接收响应, 同时defer关闭响应体
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// 读取课表
	ioRead, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "gbk") || strings.Contains(resp.Header.Get("Content-Type"), "GBK") {
		decoder := simplifiedchinese.GBK.NewDecoder()
		utf8Bytes, _, err := transform.Bytes(decoder, ioRead)
		if err != nil {
			return nil, err
		}
		return &utf8Bytes, nil
	}

	return &ioRead, nil
}

// RevLeftChidScoreUg 获取原始的本学期html成绩表(个人成绩查询), year, term为学年和学期, 要自己去html中查看
func RevLeftChidScoreUg(cookieJar *cookiejar.Jar, year, term string) (*[]byte, error) {
	//新建一个客户端
	client := &http.Client{
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	// 构建url参数
	newQuery := url.Values{}
	newQuery.Set("moduleId", "2020")
	newQuery.Set("randomString", time.Now().Format("20060102150405")+strconv.Itoa(rand.Int()))

	// 构建消息体
	values := url.Values{}
	values.Set("year", year)
	values.Set("term", term)
	values.Set("para", "0")
	values.Set("submit", "查询")

	//新建一个客户端请求
	req, err := http.NewRequest("POST", "http://jwzx.hrbust.edu.cn/academic/manager/score/studentOwnScore.do?"+newQuery.Encode(), strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", setting.UserAgent)
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//发送请求, 并接收响应, 同时defer关闭响应体
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// 读取课表
	ioRead, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &ioRead, nil
}
