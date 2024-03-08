// baidu的版本

// Package identimage 包用于通过ocr来识别验证码
package identimage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"eduData/setting"
)

// OCRResult 结构体表示整个OCR的百度的数字识别结果
type OCRResult struct {
	WordsResult []struct {
		Words    string `json:"words"`
		Location struct {
			Top    int `json:"top"`
			Left   int `json:"left"`
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"location"`
	} `json:"words_result"`
	WordsResultNum int   `json:"words_result_num"`
	LogID          int64 `json:"log_id"`
}

// NumberIdentify 传入图片base64编码, 并返回百度识别后的数据/*
func NumberIdentify(base64Image *string) (string, error) {
	requestUrl := setting.RequestUrl
	accessToken := setting.BaiduAccessToken

	client := &http.Client{} //构建http客户端实例

	values := url.Values{
		"image": {*base64Image},
	} //传入消息体
	req, err := http.NewRequest("POST", requestUrl+"?access_token="+accessToken, strings.NewReader(values.Encode()))
	if err != nil {
		return "0", errors.New("http.NewRequest()错误")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded") //增加消息头

	res, err := client.Do(req) //发送
	if err != nil {
		return "0", errors.New("client.Do()错误")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Print(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "0", errors.New("io.ReadAll()错误")
	}

	var result OCRResult
	err = json.Unmarshal(body, &result) // 解析JSON数据
	if err != nil {
		return "0", errors.New("解析JSON时发生错误")
	}

	//输出识别结果
	//fmt.Println(result)

	return result.WordsResult[0].Words, nil

}
