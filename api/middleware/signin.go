package middleware

import (
	"eduData/api/pub"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Signin 中间件登陆函数
func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		le, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
		if err != nil {
			fmt.Printf("获取失败")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": c.Error(err).Error(),
			})
			return
		}

		// 判断是哪个学校的用户来登陆
		signinCookieJar, err := pub.JudgeSchoolSignIn(loginForm)
		if err != nil {
			le.Errorf("学校登陆不成功 %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": c.Error(err).Error(),
			})
			return
		}

		// 为了后续能够获取页面, cookie加入context中
		c.Set("cookie", signinCookieJar)

		c.Next()
	}
}
