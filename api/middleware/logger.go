package middleware

import (
	"fmt"
	"net/http"
	"time"

	"eduData/domain"
	"eduData/pub"
	"eduData/school"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
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

		if param.Keys["SchoolObj"] == nil {
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-6s %s|%#-17v|%#-8v\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.ErrorMessage,
			)
		} else {
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-6s %s|%#-17v|%#-8v %#-2v %#-12v %#v %s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,

				param.Keys["SchoolObj"].(school.School).SchoolName(),
				param.Keys["SchoolObj"].(school.School).StuType(),
				param.Keys["SchoolObj"].(school.School).StuID(),
				param.Keys["SchoolObj"].(school.School).PassWd(),
				param.ErrorMessage,
			)
		}
	})
}

func CreatSchoolObject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginForm domain.LoginForm
		if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.Response{
				Status: domain.FAIL,
				Msg:    "表单格式错误, 请重新登陆 (游客模式请忽略)",
			})
			return
		}

		le := log.WithFields(log.Fields{
			"username":    loginForm.Username,
			"password":    loginForm.Password,
			"school":      loginForm.School,
			"studentType": loginForm.StudentType,
		})

		s, err := pub.NewSchoolSwitch(loginForm)
		if err != nil {
			le.Errorf("creat school obj fail: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, domain.Response{
				Status: domain.FAIL,
				Msg:    c.Error(err).Error(),
			})
			return
		}

		c.Set("SchoolObj", s)
		c.Set("logerEntry", le)
		c.Next()
	}
}
