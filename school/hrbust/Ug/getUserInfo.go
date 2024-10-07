package hrbustUg

import (
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"

	"eduData/bootstrap"
)

const SHOWHEADER = "/showPersonalInfo.do?"

func GetUserInfo(cookieJar *cookiejar.Jar) (*[]byte, error) {
	// 从setting中获取UserAgent
	userAgent := bootstrap.C.UserAgent

	// 新建一个客户端
	client := &http.Client{
		// 禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	// 构建url参数
	newQuery := url.Values{}
	newQuery.Set("randomString", time.Now().Format("20060102150405")+strconv.Itoa(rand.Int()))
	req, err := http.NewRequest(http.MethodGet, ROOT+ACADEMIC+SHOWHEADER+newQuery.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
