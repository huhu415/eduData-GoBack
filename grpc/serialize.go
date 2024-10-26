package grpc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/sirupsen/logrus"
)

const targetURL = "https://jwzx.hrbust.edu.cn/academic"

// 序列化 cookieJar
func SerializeCookieJar(jar *cookiejar.Jar) ([]byte, error) {
	ur, _ := url.Parse(targetURL)
	cookies := jar.Cookies(ur)
	logrus.Debugf("cookies: %v", cookies)

	for _, cookie := range cookies {
		logrus.Debug(cookie.Name)
		if cookie.Name == "JSESSIONID" {
			logrus.Debugf("cookies: %v", *cookie)
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if err := enc.Encode(*cookie); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
	}
	return nil, errors.New("JSESSIONID not found")
}

// DeserializeCookieJar 反序列化 cookieJar
func DeserializeCookieJar(data []byte) (*cookiejar.Jar, error) {
	// 创建一个新的 cookiejar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("DeserializeCookieJar data: %v", data)

	// 解码数据
	var cookie http.Cookie
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(&cookie); err != nil {
		return nil, err
	}

	logrus.Debugf("DeserializeCookieJar cookies: %v", cookie)

	// 将解码后的 cookie 添加到 cookiejar 中
	ur, _ := url.Parse(targetURL)
	jar.SetCookies(ur, []*http.Cookie{&cookie})

	return jar, nil
}
