package models

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"eduData/bootstrap"
)

var db *gorm.DB

// NewDatabase 新建数据库连接
func NewDatabase() {
	var err error
	dsn := bootstrap.C.PgConfig
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
	if err != nil {
		println("database : gorm.Open(), database connect error")
	}

	sqlDB, err := db.DB()
	if err != nil {
		println("database : db.DB(), database connect error")
	}

	// Migrate the schema, 创建表用的, 就用一次就完事了
	if err = db.AutoMigrate(&Course{}, &CourseGrades{}, &TimeTable{}); err != nil {
		fmt.Println("database : db.AutoMigrate(), database connect error")
		panic(err)
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(20)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(200)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("database Connected!")
}

// CloseDatabase 断开数据库连接
func CloseDatabase() {
	sqlDB, err := db.DB()
	if err != nil {
		println("database : db.DB(), database connect error")
	}
	sqlDB.Close()
	fmt.Println("database Closed!")
}
