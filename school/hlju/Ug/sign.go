package hljuUg

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"eduData/bootstrap"

	"github.com/sirupsen/logrus"
)

const (
	baseURL   = "http://xsxk.hlju.edu.cn"
	loginForm = baseURL + "/authentication/form"
)

func Signin(userName, passWord string) (*cookiejar.Jar, error) {
	userAgent := bootstrap.C.UserAgent

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Jar: cookieJar}
	defer client.CloseIdleConnections()

	// Step 1: GET 登录页，服务器下发初始 JSESSIONID + INCO cookie
	req, err := http.NewRequest(http.MethodGet, loginForm, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// Step 2: POST 用户名密码，imageCode 为空（无验证码）
	data := url.Values{}
	data.Set("username", userName)
	data.Set("password", passWord)
	data.Set("imageCode", "")

	req, err = http.NewRequest(http.MethodPost, loginForm, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+"/session/invalid")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// client 自动跟随 302，resp.Request.URL 是最终落地 URL
	finalURL := resp.Request.URL.String()
	logrus.Debugf("登录后最终URL: %s", finalURL)

	if strings.Contains(finalURL, "/authentication/form") || strings.Contains(finalURL, "session/invalid") {
		return nil, errors.New("用户名或密码错误")
	}

	logrus.Infof("登录成功，当前页面: %s", finalURL)
	return cookieJar, nil
}
