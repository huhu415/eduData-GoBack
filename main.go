package main

import (
	"eduData/bootstrap"
	"eduData/database"
	"eduData/router"
)

// todo 内存泄漏检测, 性能测试
func main() {
	bootstrap.Loadconfig()

	database.NewDatabase()
	defer database.CloseDatabase()

	router.InitRouter()
}
