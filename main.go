package main

import (
	"eduData/database"
	"eduData/router"
)

func main() {
	// todo 内存泄漏检测, 性能测试
	database.NewDatabase()
	router.InitRouter()
}
