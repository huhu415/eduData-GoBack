package main

import (
	"eduData/api/route"
	"eduData/bootstrap"
	"eduData/models"
	"github.com/sirupsen/logrus"
)

// todo 内存泄漏检测, 性能测试
func main() {
	bootstrap.Loadconfig()

	db, err := models.NewDatabase()
	if err != nil {
		logrus.Fatalf("database connect error: %v", err)
	}
	defer models.CloseDatabase(db)

	route.Setup(db)
}
