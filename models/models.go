package models

import (
	"time"

	"eduData/bootstrap"
	"eduData/repository"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase 新建数据库连接
func NewDatabase() (*gorm.DB, error) {
	dsn := bootstrap.C.PgConfig
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // precompile SQL
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Migrate the schema, 创建表用的, 就用一次就完事了
	if err = db.AutoMigrate(
		&repository.Course{},
		&repository.CourseGrades{},
		&repository.TimeTable{},
		&repository.StuInfo{}); err != nil {
		return nil, err
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(20)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(200)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Infof("\033[1;32m database Connected success! \033[0m")
	return db, nil
}

// CloseDatabase 断开数据库连接
func CloseDatabase(db *gorm.DB) {
	s, err := db.DB()
	if err != nil {
		log.Error("database : db.DB(), database connect error")
		return
	}
	if err = s.Close(); err != nil {
		log.Error("database : s.Close(), database connect error")
		return
	}
	log.Info("database Closed success!")
}
