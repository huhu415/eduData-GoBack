package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"eduData/domain"
	"eduData/pub"
	"eduData/school"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

// Prometheus 指标定义
var (
	// HTTP 请求耗时直方图
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP请求耗时分布",
			Buckets: prometheus.DefBuckets, // 默认桶: [.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10]
		},
		[]string{"method", "path", "status_code", "school", "student_type"},
	)

	// HTTP 请求总数计数器
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP请求总数",
		},
		[]string{"method", "path", "status_code", "school", "student_type"},
	)
)

// PrometheusLogger Prometheus 监控中间件
func PrometheusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Seconds()
		statusCode := strconv.Itoa(c.Writer.Status())

		// 记录指标
		school := c.Keys["SchoolObj"].(school.School)
		labels := prometheus.Labels{
			"method":       method,
			"path":         path,
			"status_code":  statusCode,
			"school":       string(school.SchoolName()),
			"student_type": strconv.Itoa(int(school.StuType())),
		}

		httpRequestDuration.With(labels).Observe(duration)
		httpRequestsTotal.With(labels).Inc()
	}
}

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
				Msg:    "表单格式错误, 请重新登陆 (游客模式请忽略)" + err.Error(),
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
