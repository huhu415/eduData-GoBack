package middleware

import (
	"eduData/api/app"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"

	"eduData/domain"
)

// Signin 中间件登陆函数
func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginForm domain.LoginForm
		if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
			_ = c.Error(errors.New("middleware.Signin()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "表单格式错误,重新登陆后重新提交",
			})
		}
		// 判断是哪个学校的用户来登陆
		signinCookieJar, err := app.JudgeSchoolSignIn(loginForm)
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": err.Error(),
			})
		}

		// 为了后续能够获取页面, cookie加入context中
		c.Set("cookie", signinCookieJar)

		c.Next()
	}
}
