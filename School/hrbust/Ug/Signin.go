package hrbustUg

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"encoding/base64"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"eduData/bootstrap"
	ident "eduData/identimage"

	"github.com/sirupsen/logrus"
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

// Signin 登陆哈理工本科生管理系统, 会拿到认证后可以使用的coookie
func Signin(USERNAME, PASSWORD string) (*cookiejar.Jar, error) {
	// 从setting中获取UserAgent
	var userAgent = bootstrap.C.UserAgent

	// 新建一个cookieJar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	//新建一个客户端
	client := &http.Client{
		// 禁止重定向
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	//第一次请求(GET方法)
	//请求首页, 从set-cookie中接收cookie
	req, err := http.NewRequest(http.MethodGet, ROOT+ACADEMIC+INDEXUG, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var OrcCodeSuccess string
	// 重试检查验证码 2 次
	for i := 0; i < 2; i++ {
		var errVer error
		var randomNum string

		// 如果不是第一次请求, 生成随机数来请求验证码图片
		if i != 0 {
			// 生成0到1之间的随机浮点数
			randomNum = fmt.Sprintf("%f", rand.Float64())
		}

		//第二次请求验证码图片
		req, errVer = http.NewRequest(http.MethodGet, ROOT+ACADEMIC+ValidateCodeSrcUg+randomNum, nil)
		if errVer != nil {
			return nil, errVer
		}
		req.Header.Add("User-Agent", userAgent)
		resp, errVer = client.Do(req)
		if errVer != nil {
			return nil, errVer
		}

		//读取回应的图片base64编码, 并识别图片
		imageBytes, errVer := io.ReadAll(resp.Body)
		if errVer != nil {
			return nil, errVer
		}
		base64String := base64.StdEncoding.EncodeToString(imageBytes)
		OrcCode, errVer := ident.CommonVerify(&base64String)
		if errVer != nil {
			return nil, errVer
		}

		//关闭图片回应体
		errVer = resp.Body.Close()
		if errVer != nil {
			return nil, errVer
		}

		//第三次请求 检查验证码是否正确
		//http://jwzx.hrbust.edu.cn/academic/checkCaptcha.do;jsessionid=40EAEB8FB62DF19770458A7C87E62C75.TA1?captchaCode=8226
		values := url.Values{}
		values.Set("captchaCode", OrcCode)
		req, errVer = http.NewRequest(http.MethodPost, ROOT+ACADEMIC+CheckCodeSrc+"captchaCode="+OrcCode, strings.NewReader(values.Encode()))
		if errVer != nil {
			return nil, errVer
		}
		req.Header.Add("User-Agent", userAgent)
		resp, errVer = client.Do(req)
		if errVer != nil {
			return nil, errVer
		}

		//读取回应的验证码识别结果
		respBody, errVer := io.ReadAll(resp.Body)
		if errVer != nil {
			return nil, errVer
		}

		//关闭验证码检查回应体
		errVer = resp.Body.Close()
		if errVer != nil {
			return nil, errVer
		}

		//如果验证码识别结果是4位, 说明回应结果为true, 跳出循环
		if len(respBody) == 4 {
			OrcCodeSuccess = OrcCode
			break
		}
	}

	// 第四次真正的登录, 带有账户密码的请求, 回来302重定向
	//http://jwzx.hrbust.edu.cn/academic/j_acegi_security_check;jsessionid=40EAEB8FB62DF19770458A7C87E62C75.TA1
	if OrcCodeSuccess != "" {
		values := url.Values{}
		values.Set("j_password", PASSWORD)
		values.Set("j_captcha", OrcCodeSuccess)
		// shit, 写成j_captch了, 少个a
		values.Set("j_username", USERNAME)
		req, err = http.NewRequest(http.MethodPost, ROOT+ACADEMIC+CHECKSECURITY, strings.NewReader(values.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		logrus.Debugf("resp.Header: %v", resp.Header)
		if strings.Contains(resp.Header.Get("Location"), "/academic/index_new.jsp") {
			//CookieParts := strings.Split(resp.Header["Set-Cookie"][0], ";")
			//CookieParts[0]就是最后登陆所要的cookie, 但既然有cookiejar, 就是用cookiejar来充当cookie
			return cookieJar, nil
		} else {
			// todo 可能要再请求过去看一下到底是怎么回事, 不过感觉暂时不用. 一般都是密码错了
			return nil, errors.New("登陆失败, 检查账号密码")
		}
	}

	return nil, errors.New("请再试一次, 自动识别验证码失败")
}
