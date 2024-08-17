// 码云的版本

// Package identimage 包用于通过ocr来识别验证码
package identimage

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"encoding/json"
	"net/http"

	"eduData/bootstrap"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Msg  string       `json:"msg"`
	Code int          `json:"code"`
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Code       int     `json:"code"`
	Data       string  `json:"data"`
	Time       float64 `json:"time"`
	Externel   int     `json:"externel"`
	UniqueCode string  `json:"unique_code"`
}

func CommonVerify(image *string) (string, error) {
	CustomUrl := bootstrap.C.JfymRequestUrl
	Token := bootstrap.C.JfymToken

	//通用数英1~6位plus 10103
	config := map[string]any{}
	config["image"] = *image
	config["type"] = "10103"
	config["token"] = Token
	configData, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(configData)
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	logrus.Debugf("内容: %v,云码平台识别结果: %v", config, string(data))

	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", fmt.Errorf("自动识别验证码错误: %v", err)
	}

	// https://www.jfbym.com/demo.html
	switch response.Code {
	case 10000:
		return response.Data.Data, nil
	case 10001:
		return "", errors.New("参数错误")
	case 10002:
		return "", errors.New("余额不足")
	case 10003:
		return "", errors.New("无此访问权限")
	case 10004:
		return "", errors.New("无此验证类型")
	case 10005:
		return "", errors.New("网络拥塞")
	case 10006:
		return "", errors.New("数据包过载")
	case 10007:
		return "", errors.New("服务繁忙")
	case 10008:
		return "", errors.New("网络错误，请稍后重试")
	case 10009:
		return "", errors.New("结果准备中，请稍后再试")
	case 10010:
		return "", errors.New("请求结束")
	}
	return "", errors.New("云码平台识别未成功" + response.Msg + string(rune(response.Code)))
}
