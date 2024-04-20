package main

import (
	"eduData/api"
	"eduData/bootstrap"
	"eduData/models"
)

// todo 内存泄漏检测, 性能测试
func main() {
	bootstrap.Loadconfig()

	models.NewDatabase()
	defer models.CloseDatabase()

	api.InitRouterRunServer()
}
