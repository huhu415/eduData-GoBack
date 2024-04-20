package middleware

import (
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	hrbustPg "eduData/School/hrbust/Pg"
	hrbustUg "eduData/School/hrbust/Ug"
	neauUg "eduData/School/neau/Ug"
)

// JudgeSchoolSignIn 判断是哪个学校的用户来登陆
func judgeSchoolSignIn(loginForm LoginForm) (*cookiejar.Jar, error) {
	switch loginForm.School {
	// 哈理工
	case "hrbust":
		switch loginForm.StudentType {
		case 1:
			return hrbustUg.Signin(loginForm.Username, loginForm.Password)
		case 2:
			return hrbustPg.Signin(loginForm.Username, loginForm.Password)
		}
	// 东北农业大学
	case "neau":
		switch loginForm.StudentType {
		case 1:
			return neauUg.Signin(loginForm.Username, loginForm.Password)
		case 2:
			return nil, errors.New(loginForm.School + "研究生登陆功能还未开发")
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
		var loginForm LoginForm
		if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
			_ = c.Error(errors.New("middleware.Signin()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "表单格式错误,重新登陆后重新提交",
			})
			c.Abort()
		}
		// 判断是哪个学校的用户来登陆
		signinCookieJar, err := judgeSchoolSignIn(loginForm)
		if err != nil {
			_ = c.Error(errors.New("middleware.signin()JudgeSchoolSignIn():" + err.Error())).SetType(gin.ErrorTypePrivate)

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
