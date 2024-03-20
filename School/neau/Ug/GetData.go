package neauUg

import (
	"eduData/setting"
	"io"
	"net/http"
	"net/http/cookiejar"
)

// GetData 获取课表json data是2023-2024-2-1(23年第二学期) 2024-2025-1-1(24年第一学期)
func GetData(cookieJar *cookiejar.Jar, data string) (*[]byte, error) {
	//新建一个客户端, 运行重定向, 设置cookie
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 允许重定向
			return nil
		},
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	// 发送请求, 获取课表, 收到json格式
	request, err := http.NewRequest("POST", "https://zhjwxs-443.webvpn.neau.edu.cn/student/courseSelect/thisSemesterCurriculum/ajaxStudentSchedule/callback", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", setting.UserAgent)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("planCode", data)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	readAll, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &readAll, nil
}
