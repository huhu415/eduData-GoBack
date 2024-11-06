package hljuUg

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"eduData/bootstrap"
)

const (
	kbjg            = "http://xsxk.hlju.edu.cn/component/queryKbjg"
	queryxszykbzong = "http://xsxk.hlju.edu.cn/xszykb/queryxszykbzong"
)

// for example: now is 2024-11-6 the xn is 2024-2025 xq is 1, xn is 2024
// if now is 2025-3-6 the xn is 2024-2025 xq is 2, xn is 2024
func GetData(cookie *cookiejar.Jar, xn int, xq int) (*[]byte, error) {
	userAgent := bootstrap.C.UserAgent
	client := &http.Client{Jar: cookie}

	values := url.Values{}
	values.Set("xn", fmt.Sprintf("%d-%d", xn, xn+1))
	values.Set("xq", fmt.Sprintf("%d", xq))
	data := strings.NewReader(values.Encode())

	// var data = strings.NewReader(`xn=2024-2025&xq=1&pylx=1`)
	req, err := http.NewRequest(http.MethodPost, queryxszykbzong, data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "http://xsxk.hlju.edu.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://xsxk.hlju.edu.cn/authentication/main")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return &bodyText, nil
}
