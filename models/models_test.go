package models

import (
	"eduData/bootstrap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	bootstrap.InitLog()
	bootstrap.Loadconfig()
	NewDatabase()

	NewDatabase()
	defer CloseDatabase()

	t.Run("添加课程", func(t *testing.T) {

	})

	t.Run("查询课程", func(t *testing.T) {
		courses := CourseByWeekUsername(1, "1234567", "hrbust")
		t.Log(courses)
	})

	t.Run("查询成绩", func(t *testing.T) {
		gpa1, gpa2 := WeightedAverage("2204010417", "hrbust", 1)
		t.Log(gpa1, gpa2)
	})

	t.Run("更新/添加个人信息", func(t *testing.T) {
		stu := &StuInfo{
			StuID:   "2306070112",
			School:  "hrbust",
			StuType: 1,
		}
		err := stu.CreatAndUpdataStuInfo()
		assert.Nil(t, err, "更新个人信息失败")
		t.Log(stu)
	})
}
