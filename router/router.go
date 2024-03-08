package router

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"eduData/app"
	"eduData/middleware"
	"eduData/setting"
)

func InitRouter() {
	// 发布模式, 删了就是debug模式
	gin.SetMode(gin.ReleaseMode)

	// 日志记录
	f, err := os.OpenFile("crawler.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			fmt.Println("关闭文件失败" + err.Error())
		}
	}(f)
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// 路由初始化, 1.日志 2.恢复 3.检查表单完成性
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery(), middleware.LoggerRecordForm())

	// 路由分组
	si := r.Group("/")
	si.Use(middleware.Signin())
	{
		si.POST("/signin", app.Signin)
		si.POST("/updata", app.UpdataDB)
	}
	auth := r.Group("/getweekcoure")
	auth.Use(middleware.RequireAuthJwt())
	//auth.Use()
	{
		auth.POST("/:week", app.GetWeekCoure)
	}
	// todo 增加内部测试功能, 内部账号密码登录, 发送一个html, 返回解析好的课表, 可以用周数的那个函数来调用

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", setting.HttpPort),
		Handler: r,
	}

	// https模式
	//go func() {
	//	if err = srv.ListenAndServeTLS("ssl/zzyan.com_ssl.crt", "ssl/zzyan.com_ssl.key"); err != nil || errors.Is(http.ErrServerClosed, err) {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()

	// 普通模式
	go func() {
		if err = srv.ListenAndServe(); err != nil || errors.Is(http.ErrServerClosed, err) {
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
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
