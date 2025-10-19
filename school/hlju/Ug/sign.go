package hljuUg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"eduData/bootstrap"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

const (
	ROOT string = "https://authserver.hlju.edu.cn/authserver/login"
	CAS  string = "http://xsxk.hlju.edu.cn/cas"
)

func Signin(userName, passWord string) (*cookiejar.Jar, error) {
	userAgent := bootstrap.C.UserAgent

	// 新建一个cookieJar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	// 新建一个客户端, 自动重定向
	client := &http.Client{
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	// 第一次请求(GET方法)
	// 请求首页, 从set-cookie中接收cookie
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

	// jsessionid=aOSyo_sGDu_6cRt3n7jqhIJ-JH2QVnFu_0CWiSfIpopk3-T-a3mY!1507668027
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
		// 到这里已经登录成功
		time.Sleep(1 * time.Second)
		GetCas(client)
		return cookieJar, nil
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

// GET http://xsxk.hlju.edu.cn/cas for get course 'queryKbjg'
func GetCas(client *http.Client) error {
	req, err := http.NewRequest(http.MethodGet, CAS, nil)
	if err != nil {
		logrus.Error(err)
		return err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "zh,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("DNT", "1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://xsxk.hlju.edu.cn/")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", bootstrap.C.UserAgent)
	// 允许重定向
	client.CheckRedirect = nil

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()
	return nil
}
