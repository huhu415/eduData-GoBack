package main

import (
	"eduData/api/route"
	"eduData/bootstrap"
	"eduData/models"

	"github.com/sirupsen/logrus"
)

func main() {
	bootstrap.Loadconfig()

	db, err := models.NewDatabase()
	if err != nil {
		logrus.Fatalf("Fatal!! database connect error: %v", err)
	}
	defer models.CloseDatabase(db)

	route.SetupAndRun(db)
}
