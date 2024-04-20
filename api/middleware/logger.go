package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"eduData/domain"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v %#v %#v %#v %#v \n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.Keys["username"],
			param.Keys["password"],
			param.Keys["school"],
			param.Keys["studentType"],
			param.ErrorMessage,
		)
	})
}

// LoggerRecordForm 记录提交的表单中的内容
func LoggerRecordForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginForm domain.LoginForm
		if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
			_ = c.Error(errors.New("middleware.LoggerRecordForm()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "表单格式错误,重新登陆后重新提交",
			})
		}

		c.Set("username", loginForm.Username)
		c.Set("password", loginForm.Password)
		c.Set("school", loginForm.School)
		c.Set("studentType", loginForm.StudentType)
		c.Next()
	}

}
