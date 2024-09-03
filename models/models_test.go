package models

import (
	"eduData/bootstrap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
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

	t.Run("添加时间表", func(t *testing.T) {
		times := make([]TimeTable, 0)
		for i, v := range courses {
			times = append(times, TimeTable{
				School:    "hlju",
				Sort:      uint(i + 1),
				StartTime: v.startTime,
				EndTime:   v.endTime,
			})
		}
		assert.Nil(t, AddTimeTable(&times), "添加时间表失败")
	})
}

type course struct {
	startTime string
	endTime   string
}

// hlju时间表
var courses = []course{
	{"08:00", "08:45"},
	{"08:50", "09:35"},
	{"10:00", "10:45"},
	{"10:50", "11:35"},
	{"13:30", "14:15"},
	{"14:20", "15:05"},
	{"15:30", "16:15"},
	{"16:20", "17:05"},
	{"18:30", "19:15"},
	{"19:20", "20:05"},
	{"20:10", "20:55"},
}
