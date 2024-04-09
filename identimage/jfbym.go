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
	config := map[string]interface{}{}
	config["image"] = *image
	config["type"] = "10103"
	config["token"] = Token
	configData, err := json.Marshal(config)
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer(configData)
	resp, err := http.Post(CustomUrl, "application/json;charset=utf-8", body)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var response Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}
	// 10000 为成功, https://www.jfbym.com/test/52.html详情
	if response.Code == 10000 {
		return response.Data.Data, nil
	}
	return "", errors.New("云码平台识别未成功" + response.Msg + string(rune(response.Code)))
}
