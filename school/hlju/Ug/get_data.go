package hljuUg

import (
	"eduData/bootstrap"
	"io"
	"net/http"
	"net/http/cookiejar"
)

func GetData(cookie *cookiejar.Jar) (*[]byte, error) {
	var userAgent = bootstrap.C.UserAgent

	//新建一个客户端
	client := &http.Client{
		Jar: cookie,
	}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest(http.MethodGet, "http://ssfw1.hlju.edu.cn/ssfw/xkgl/xkjgcx.do", nil)
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
