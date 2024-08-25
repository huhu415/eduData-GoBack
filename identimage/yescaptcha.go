package identimage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const IMAGE_TO_TEXT_TASK_TEST = "ImageToTextTaskTest"

type YescaptchaOcr struct {
	ClientUrl   string
	ClientToken string
}

func NewYescaptchaOcr(url, token string) IdentImage {
	return &YescaptchaOcr{
		ClientUrl:   url,
		ClientToken: token,
	}
}

func (y *YescaptchaOcr) Identify(base64Image *string) (string, error) {
	// 构建请求数据
	reqData := yescaptchaRequestData{
		ClientKey: y.ClientToken,
		Task: task{
			Type: IMAGE_TO_TEXT_TASK_TEST,
			Body: *base64Image,
		},
	}

	// 将请求数据编码为JSON
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		fmt.Println("Error encoding request data:", err)
		return "", err
	}

	// 发送HTTP POST请求
	resp, err := http.Post(y.ClientUrl, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return "", err
	}

	// 解析JSON响应
	var respData yescaptchaResponseData
	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Println("Error decoding response data:", err)
		return "", err
	}

	// 输出响应结果
	if respData.ErrorId == 0 && respData.Status == "ready" {
		fmt.Println("YescaptchaOcr Task successful")
		return respData.Solution.Text, nil
	} else {
		return "", fmt.Errorf("task failed: %s (%s)", respData.ErrorDescription, respData.ErrorCode)
	}
}

type task struct {
	Type string `json:"type"`
	Body string `json:"body"`
}

type yescaptchaRequestData struct {
	ClientKey string `json:"clientKey"`
	Task      task   `json:"task"`
}
type yescaptchaResponseData struct {
	ErrorId          int    `json:"errorId"`
	ErrorCode        string `json:"errorCode"`
	ErrorDescription string `json:"errorDescription"`
	Status           string `json:"status"`
	Solution         struct {
		Text string `json:"text"`
	} `json:"solution"`
	TaskId string `json:"taskId"`
}
