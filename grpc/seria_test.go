package grpc

import (
	"net/http"
	"net/http/cookiejar"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	ROOT     string = "http://jwzx.hrbust.edu.cn"
	ACADEMIC string = "/academic"
	// INDEXUG 哈理工研究生首页地址
	INDEXUG string = "/index.jsp"
	// ValidateCodeSrcUg 请求验证码图片地址
	ValidateCodeSrcUg string = "/getCaptcha.do?"
	// CheckCodeSrc 检查验证码地址
	CheckCodeSrc string = "/checkCaptcha.do?"
	// CHECKSECURITY 最终登陆地址
	CHECKSECURITY string = "/j_acegi_security_check"
)

func TestSerializeDeserializeCookieJar(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	assert := assert.New(t)
	// 创建一个新的 cookiejar
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("Failed to create cookiejar: %v", err)
	}

	// 新建一个客户端
	client := &http.Client{
		// 禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}
	defer client.CloseIdleConnections()

	// 第一次请求(GET方法)
	// 请求首页, 从set-cookie中接收cookie
	req, err := http.NewRequest(http.MethodGet, ROOT+ACADEMIC+INDEXUG, nil)
	assert.Nil(err)

	req.Header.Add("User-Agent", "FFFFFFFFFFFFFFFF")
	resp, err := client.Do(req)
	assert.Nil(err)

	defer resp.Body.Close()

	// 序列化 cookiejar
	serializedData, err := SerializeCookieJar(jar)
	if err != nil {
		t.Fatalf("Failed to serialize cookiejar: %v", err)
	}

	// 反序列化 cookiejar
	_, err = DeserializeCookieJar(serializedData)
	if err != nil {
		t.Fatalf("Failed to deserialize cookiejar: %v", err)
	}
}
