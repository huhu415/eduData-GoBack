package hljuUg

import (
	"bytes"
	"encoding/json"
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

type courseInfo struct {
	Xn       *string `json:"xn"`       // 学年 (Academic Year)，可能为空
	Xq       *string `json:"xq"`       // 学期 (Semester)，可能为空
	Kcmc     *string `json:"kcmc"`     // 课程名称 (Course Name)，可能为空
	Cxbj     string  `json:"cxbj"`     // 重修标记 (Retake Flag)，"-1" 表示未定义或不重修
	Pylx     string  `json:"pylx"`     // 培养类型 (Training Type)，例如 "1" 代表某种特定的培养类型
	Current  int     `json:"current"`  // 当前页码 (Current Page)，用于分页的当前页
	PageSize int     `json:"pageSize"` // 每页大小 (Page Size)，用于分页的每页记录数
	Sffx     *string `json:"sffx"`     // 是否分析 (Analysis Flag)，可能为空
}

func GetScore(cookie *cookiejar.Jar) (*[]byte, error) {
	userAgent := bootstrap.C.UserAgent
	client := &http.Client{Jar: cookie}

	c := courseInfo{
		Cxbj:     "-1",
		Pylx:     "1",
		Current:  1,
		PageSize: 200,
	}
	jsonC, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, "http://xsxk.hlju.edu.cn/cjgl/grcjcx/grcjcx", bytes.NewReader(jsonC))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("DNT", "1")
	req.Header.Set("Origin", "http://xsxk.hlju.edu.cn")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Proxy-Connection", "keep-alive")
	req.Header.Set("Referer", "http://xsxk.hlju.edu.cn/cjgl/grcjcx/go/1")
	req.Header.Set("RoleCode", "01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", userAgent)

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
