package hrbustPg

import (
	"errors"
	"fmt"
	"io"
	"net/http/cookiejar"
	"strings"

	"encoding/base64"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"

	"eduData/bootstrap"
	ident "eduData/identimage"
)

const (
	// INDEXPG 哈理工研究生首页地址
	INDEXPG string = "http://yjs.hrbust.edu.cn/"
)

// Signin 登陆哈理工研究生管理系统, 会拿到认证后可以使用的coookie
func Signin(USERNAME, PASSWORD string) (*cookiejar.Jar, error) {
	// 从setting中获取UserAgent
	var userAgent = bootstrap.C.UserAgent

	//新建一个客户端
	client := &http.Client{}
	defer client.CloseIdleConnections()

	//第一次请求(GET方法)
	//请求首页, 接收带有__VIEWSTATE, __EVENTVALIDATION, ValidateCode链接的响应体,
	//同时__VIEWSTATE里面也有验证码ValidateCode链接, 两个途径都可以拿到ValidateCode链接
	req, err := http.NewRequest("GET", INDEXPG, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	//从得到的*http.Response中的body得到__VIEWSTATE和__EVENTVALIDATION, 并且通过__VIEWSTATE得到验证码ValidateCode的请求地址
	VIEWSTATEvalue, EVENTVALIDATIONvalue, ValidateCodeSrc, err := findVEI(resp)
	if err != nil {
		return nil, err
	}

	//第二次请求(GET方法)
	//请求验证码, 并接收图片和cookie
	req, err = http.NewRequest("GET", INDEXPG+ValidateCodeSrc, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", userAgent)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("defer func(Body io.ReadCloser) error: %v\n", err)
		}
	}(resp.Body)

	//从*http.Response中得到cookie和验证码图片, 并将验证码识别成字符串类型的数字
	IdentifyNum, cookie, err := findValidateCodeCookie(resp)
	if err != nil {
		return nil, err
	}

	//第三次请求(POST方法)
	//登陆过程的最后一步, 就是把ValidateCode, UserName, PassWord, __VIEWSTATE, __EVENTVALIDATION
	//这几个会改变的值一起post给哈理工研究生管理系统
	values := url.Values{}
	values.Set("ScriptManager1", "UpdatePanel2|btLogin")
	values.Set("__EVENTTARGET", "btLogin")
	values.Set("__EVENTARGUMENT", "")
	values.Set("__LASTFOCUS", "")
	values.Set("__VIEWSTATE", VIEWSTATEvalue)
	values.Set("__EVENTVALIDATION", EVENTVALIDATIONvalue)
	values.Set("UserName", USERNAME)
	values.Set("PassWord", PASSWORD)
	values.Set("ValidateCode", IdentifyNum)
	values.Set("drpLoginType", "1")
	values.Set("__ASYNCPOST", "true")
	//构建请求体
	req, _ = http.NewRequest("POST", INDEXPG, strings.NewReader(values.Encode()))
	req.Header.Add("cookie", cookie)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			fmt.Printf("defer func(Body io.ReadCloser) error: %v\n", err)
		}
	}(resp.Body)

	ResBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll()错误: %v", err)
	}

	//判断登陆成功与否
	if (len(ResBody) < 200 && len(ResBody) > 50) && (strings.Contains(string(ResBody), USERNAME)) {
		//登陆成功
		//fmt.Println("sign.SingInPg()登陆成功,消息体:" + string(ResBody))
		//fmt.Println("登陆成功, cookie为: ", cookie)
		// cookie的样子是ASP.NET_SessionId=cjc3oisjjnenugay4cy0xa4m(; path=/; HttpOnly)这个被findValidateCodeCookie()函数去掉了
		CookieParts := strings.Split(cookie, "=")
		// 新建一个cookieJar
		cookieJar, err := cookiejar.New(nil)
		if err != nil {
			return nil, err
		}
		cookieURL, err := url.Parse(INDEXPG)
		if err != nil {
			return nil, err
		}
		cookieJar.SetCookies(cookieURL, []*http.Cookie{
			{
				Name:  CookieParts[0],
				Value: CookieParts[1],
			},
		})
		return cookieJar, nil
	} else {
		//登陆失败
		//fmt.Println("sign.SingInPg()登陆失败:", string(ResBody[i-40:]))
		switch {
		case strings.Contains(string(ResBody), "密码输入错误"):
			return nil, errors.New("密码输入错误")
		case strings.Contains(string(ResBody), "验证码错误"):
			return nil, errors.New("验证码错误")
		case strings.Contains(string(ResBody), "用户帐号不存在"):
			return nil, errors.New("用户帐号不存在")
		case strings.Contains(string(ResBody), "验证码输入错误"):
			return nil, errors.New("验证码输入错误")
		}
	}

	//判断登陆成功与否
	return nil, errors.New(string(ResBody[len(ResBody)-40:]))
}

// findValidateCodeCookie 从*http.Response中得到cookie和验证码图片,
// 将验证码识别成字符串类型的数字, 需要我自己写的eduData/identimage包*/
func findValidateCodeCookie(resp *http.Response) (IdentifyNum, cookieres string, err error) {
	cookie := resp.Header["Set-Cookie"]
	CookieParts := strings.Split(cookie[0], ";")
	cookieres = CookieParts[0]
	//从原生cookie的值解析出真正的cookie, 名为CookieParts

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", "0", errors.New("io.ReadAll()错误" + err.Error())
	}
	base64String := base64.StdEncoding.EncodeToString(imageBytes)
	base64String = strings.TrimPrefix(base64String, "data:image/gif;base64,")
	// 使用 strings 包去除可能的Base64编码时添加的前缀（如 "data:image/gif;base64,"）最后得到base64编码后的图片数据

	IdentifyNum, err = ident.NumberIdentify(&base64String)
	if err != nil {
		return "0", "0", err
	}
	//调用百度的ocr接口识别验证码并返回结果
	return IdentifyNum, cookieres, nil
}

// findVEI 从得到的*http.Response中的body分别得到__VIEWSTATE和__EVENTVALIDATION的value值,
// 并且通过__VIEWSTATE得到验证码ValidateCode的请求地址*/
func findVEI(resp *http.Response) (VIEWSTATEvalue, EVENTVALIDATIONvalue, ValidateCodeSrc string, err error) {
	read, err := io.ReadAll(resp.Body)
	if err != nil {
		return "0", "0", "0", err
	}

	// 使用 goquery 解析 HTML 文档
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(read)))
	if err != nil {
		return "0", "0", "0", err
	}

	//第一种方式
	//解析响应体, 并找到__VIEWSTATE和__EVENTVALIDATION, ValidateCode链接的值
	VIEWSTATEvalue, ok := doc.Find("#__VIEWSTATE").Attr("value")
	if !ok {
		return "0", "0", "0", errors.New("not find __VIEWSTATE")
	}
	EVENTVALIDATIONvalue, ok = doc.Find("#__EVENTVALIDATION").Attr("value")
	if !ok {
		return "0", "0", "0", errors.New("not find __EVENTVALIDATION")
	}
	ValidateCodeSrc, ok = doc.Find("#ValidateImage").Attr("src")
	if !ok {
		return "0", "0", "0", errors.New("not find ValidateImage")
	}

	/*第二种方式
	//通过__VIEWSTATE来base64解码来, 找到ValidateCode链接的值
		decodedBytes, err := base64.StdEncoding.DecodeString(VIEWSTATEvalue)
		if err != nil {
			return "0", "0", "0", err
		}
		// 进行Base64解码,解码后会有子串类似于image=1234567890dd
		ImageIndex := strings.Index(string(decodedBytes), "image=")
		if ImageIndex == -1 {
			return "0", "0", "0", err
		}
		for k := ImageIndex + 6; string(decodedBytes[k]) != "d"; k++ {
			ImageNumber += string(decodedBytes[k])
			if k > ImageIndex+50 {
				return "0", "0", "0", errors.New("没有找到image=数组后面的dd")
			}
		}
		//通过解码后的字符串, 找到image=后面的数字, 并赋值给ImageNumber, 这个数字就是验证码的图片名

		ValidateCodeSrc = "Public/ValidateCode.aspx?image=" + ImageNumber
		//拼接出验证码的请求地址
	*/
	return VIEWSTATEvalue, EVENTVALIDATIONvalue, ValidateCodeSrc, nil
}
