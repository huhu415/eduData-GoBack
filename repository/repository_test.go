package repository

import (
	"eduData/bootstrap"
	"eduData/school/pub"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository(t *testing.T) {
	bootstrap.Loadconfig()
	dsn := bootstrap.C.PgConfig
	db, err := gorm.Open(postgres.Open(dsn))
	assert.Nil(t, err, "database connect error")

	ur := NewRepository(db)

	t.Run("TestAddCourse", func(t *testing.T) {
		course, err := ur.CourseByGroup("A05250061", pub.NEAU, pub.UG)
		assert.NoError(t, err, "course by group error")
		json, err := json.Marshal(course)
		assert.NoError(t, err, "marshal error")
		fmt.Println(string(json))
	})

}
