package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"eduData/api/app"
	"eduData/api/middleware"
	"eduData/bootstrap"
)

func InitRouterRunServer() {
	log.SetReportCaller(true)
	// 发布模式, 删了就是debug模式
	gin.SetMode(gin.ReleaseMode)

	// 日志记录
	gin.ForceConsoleColor()
	f, err := os.OpenFile("eduData.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetOutput(gin.DefaultWriter)

	// 路由初始化, 1.日志 2.恢复 3.检查表单完成性
	r := gin.New()
	r.Any("/health", gin.Logger(), gin.Recovery(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "active",
		})
	})
	r.Use(middleware.Logger(), gin.Recovery(), middleware.LoggerRecordForm())

	si := r.Group("/")
	si.Use(middleware.Signin())
	{
		si.POST("/signin", app.Signin)
		si.POST("/updata", app.UpdataDB)
		si.POST("/updataGrade", app.UpdataGrade)
	}

	auth := r.Group("/")
	auth.Use(middleware.RequireAuthJwt())
	//auth.Use()
	{
		auth.POST("/getweekcoure/:week", app.GetWeekCoure)
		auth.POST("/getgrade", app.GetGrade)
		auth.POST("/getTimeTable", app.GetTimeTable)
		auth.POST("/addcoures", app.AddCoures)
	}
	// todo 增加内部测试功能, 内部账号密码登录, 发送一个html, 返回解析好的课表, 可以用周数的那个函数来调用
	// todo https://gin-gonic.com/zh-cn/docs/examples/serving-data-from-reader/

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", bootstrap.C.ListenPort),
		Handler: r,
	}
	runServer(srv)
}

func runServer(srv *http.Server) {
	go func() {
		log.Infof("Server start listening at %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil || errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Infof("Server exiting")
}
