package hljuUg

import (
	"bytes"
	"eduData/bootstrap"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const ROOT string = "https://authserver.hlju.edu.cn/authserver/login"

func Signin(userName, passWord string) (*cookiejar.Jar, error) {
	var userAgent = bootstrap.C.UserAgent

	// 新建一个cookieJar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	//新建一个客户端, 自动重定向
	client := &http.Client{
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	//第一次请求(GET方法)
	//请求首页, 从set-cookie中接收cookie
	fullURL := fmt.Sprintf("%s?", ROOT)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//jsessionid=aOSyo_sGDu_6cRt3n7jqhIJ-JH2QVnFu_0CWiSfIpopk3-T-a3mY!1507668027
	cookie := ""
	for _, c := range resp.Cookies() {
		if c.Name == "JSESSIONID" {
			logrus.Debugf("JSESSIONID: %s", c.Value)
			cookie = c.Value
			break
		}
	}

	Lt, err := getLT(resp)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("lt: %s", Lt)

	// 请求体
	data := url.Values{}
	data.Set("username", userName)
	data.Set("password", passWord)
	data.Set("btn", "")
	data.Set("lt", Lt)
	data.Set("dllt", "userNamePasswordLogin")
	data.Set("execution", "e1s1")
	data.Set("_eventId", "submit")
	data.Set("rmShown", "1")
	fullURL = fmt.Sprintf("%s;jsessionid=%s", ROOT, cookie)
	request, err := http.NewRequest(http.MethodPost, fullURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", userAgent)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// 禁止重定向, 用于检查相应头的location
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 && strings.Contains(resp.Header.Get("Location"), "index.do") {
		time.Sleep(2 * time.Second)
		req, err = http.NewRequest(http.MethodGet, "http://ssfw1.hlju.edu.cn/ssfw/j_spring_ids_security_check", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("User-Agent", userAgent)
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == 302 && strings.Contains(resp.Header.Get("Location"), "success=true") {
			logrus.Debug("success Signin")
			return cookieJar, nil
		} else {
			return nil, errors.New("登陆失败")
		}
	} else if resp.StatusCode == 200 {
		return nil, errors.New("您提供的用户名或者密码有误")
	}

	return nil, errors.New("未知错误")
}

func getLT(resp *http.Response) (string, error) {
	// 从一个html中, 找到lt
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	LTvalue, ok := doc.Find("input[name='lt']").Attr("value")
	if ok {
		return LTvalue, nil
	} else {
		return "", errors.New("lt value not found")
	}
}
