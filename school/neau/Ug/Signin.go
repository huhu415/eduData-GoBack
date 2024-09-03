package neauUg

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"net/http/cookiejar"

	"github.com/PuerkitoBio/goquery"

	"eduData/bootstrap"
)

// encryptPassword 加密密码
func encryptPassword(n, f string) (string, error) {
	ciphertext, err := getAesString(randomString(64)+n, f, randomString(16))
	if err != nil {
		return "", err
	}
	return ciphertext, nil
}

// randomString 从已知字符串中随机获取n个字符
func randomString(n int) string {
	seed := time.Now().UnixNano()
	rand.New(rand.NewSource(seed))
	aesChars := "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678"
	b := make([]byte, n)
	for i := range b {
		b[i] = aesChars[rand.Intn(len(aesChars))]
	}
	return string(b)
}

// getAesString AES加密
func getAesString(n, f, c string) (string, error) {
	// Trim spaces from the key and initialization vector
	f = string(bytes.TrimSpace([]byte(f)))
	c = string(bytes.TrimSpace([]byte(c)))

	// Parse the key and initialization vector
	key, err := aes.NewCipher([]byte(f))
	if err != nil {
		return "", err
	}
	iv := []byte(c)

	// Pad the plaintext
	padLength := aes.BlockSize - len(n)%aes.BlockSize
	pad := bytes.Repeat([]byte{byte(padLength)}, padLength)
	plaintext := append([]byte(n), pad...)

	// Encrypt the plaintext
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(key, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// Encode the ciphertext to base64 string
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)

	return ciphertextBase64, nil
}

// Signin 农大登陆, 得到cookieJar
func Signin(username, password string) (*cookiejar.Jar, error) {
	// 新建一个cookieJar
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	//新建一个客户端
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 允许重定向
			return nil
		},
		Jar: cookieJar,
	}
	defer client.CloseIdleConnections()

	// 发送请求第一次, 找到execution和pwdEncryptSalt
	request, err := http.NewRequest("GET", "https://authserver-443.webvpn.neau.edu.cn/authserver/login", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", bootstrap.C.UserAgent)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 从响应体中获取execution和pwdEncryptSalt
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	execution, ok := doc.Find("input#execution").Attr("value")
	if !ok {
		return nil, err
	}
	pwdEncryptSalt, ok := doc.Find("input#pwdEncryptSalt").Attr("value")
	if !ok {
		return nil, err
	}

	// 加密密码
	password, err = encryptPassword(password, pwdEncryptSalt)
	if err != nil {
		return nil, err
	}

	// 构建消息体
	values := url.Values{}
	values.Set("username", username)
	values.Set("password", password)
	values.Set("_eventId", "submit")
	values.Set("cllt", "userNameLogin")
	values.Set("dllt", "generalLogin")
	values.Set("execution", execution)
	// 发送请求第二次, 登陆
	request, err = http.NewRequest(http.MethodPost, "https://authserver-443.webvpn.neau.edu.cn/authserver/login", strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", bootstrap.C.UserAgent)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(all) > 10000 {
		return nil, errors.New("登陆失败, 检查账号密码")
	}

	// 发送请求第三次, 获取到主页, 因为获取主页的过程, 有各种cookie的添加
	request, err = http.NewRequest("GET", "https://zhjwxs-443.webvpn.neau.edu.cn/index", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", bootstrap.C.UserAgent)
	resp, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return cookieJar, nil
}
