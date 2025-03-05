package route

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

	"eduData/api/middleware"
	"eduData/bootstrap"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupAndRun(db *gorm.DB) {
	if viper.GetBool("debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 日志记录
	gin.ForceConsoleColor()
	f, err := os.OpenFile("eduData.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatal("Fatal!! OpenFile failed: ", err)
	}
	defer f.Close()
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	log.SetOutput(gin.DefaultWriter)

	// 路由初始化, 1.日志 2.恢复 3.检查表单完成性
	r := gin.New()
	r.GET("/health", gin.Recovery(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "active",
		})
	})
	r.Use(middleware.Logger(), gin.Recovery(), middleware.CreatSchoolObject())

	NewUpDataRouter(db, r.Group(""))

	runServer(&http.Server{
		Addr:    fmt.Sprintf(":%s", bootstrap.C.ListenPort),
		Handler: r,
	})
}

func runServer(srv *http.Server) {
	go func() {
		logrus.Infof("\u001B[1;32m Server start listening at%s! \u001B[0m", srv.Addr)
		logrus.SetReportCaller(true)
		if err := srv.ListenAndServe(); err != nil || errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("Fatal!! listen: %s\n", err)
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
	logrus.Infof("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Fatal!! Server forced to shutdown: ", err)
	}

	logrus.Infof("Server already shutdown")
}
