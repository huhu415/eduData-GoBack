package middleware

import (
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/gin-gonic/gin"

	signin "eduData/sign_in"
)

// JudgeSchoolSignIn 判断是哪个学校的用户来登陆
func judgeSchoolSignIn(c *gin.Context) (*cookiejar.Jar, error) {
	switch c.PostForm("school") {
	// 哈理工
	case "hrbust":
		switch c.PostForm("studentType") {
		case "1":
			return signin.SingInUg(c.PostForm("username"), c.PostForm("password"))
		case "2":
			return signin.SingInPg(c.PostForm("username"), c.PostForm("password"))
		}
	// 东北农业大学
	case "neau":
		switch c.PostForm("studentType") {
		case "1":
			return signin.SigninUgNEAU(c.PostForm("username"), c.PostForm("password"))
		case "2":
			return nil, errors.New(c.PostForm("school") + "研究生登陆功能还未开发")
		}
	// 其他没有适配的学校
	default:
		return nil, errors.New("不支持的学校")
	}
	return nil, errors.New("不支持的学校")
}

// Signin 中间件登陆函数
func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 判断是哪个学校的用户来登陆
		signinCookieJar, err := judgeSchoolSignIn(c)
		if err != nil {
			_ = c.Error(errors.New("middleware.signin()JudgeSchoolSignIn():" + c.PostForm("school") + c.PostForm("studentType") + err.Error())).SetType(gin.ErrorTypePrivate)

			// 获取最内层的错误给用户
			for {
				nextErr := errors.Unwrap(err)
				if nextErr == nil {
					break
				}
				err = nextErr
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
			c.Abort()
		}

		// 为了后续能够获取页面, cookie加入context中
		c.Set("cookie", signinCookieJar)

		c.Next()
	}
}
