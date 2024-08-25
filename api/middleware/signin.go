package middleware

import (
	"eduData/pub"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Signin 中间件登陆函数
func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		s, le, err := pub.GetSchoolAndLogrus(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": c.Error(err).Error(),
			})
			return
		}

		if err := s.Signin(); err != nil {
			le.Errorf("登陆失败 %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": c.Error(err).Error(),
			})
			return
		}
		c.Next()
	}
}
