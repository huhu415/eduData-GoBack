package models

import (
	"testing"

	"eduData/bootstrap"
	"eduData/repository"
	"eduData/school/pub"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	bootstrap.Loadconfig()
	db, err := NewDatabase()
	assert.Nil(t, err, "连接数据库失败")
	defer CloseDatabase(db)

	ur := repository.NewRepository(db)

	t.Run("添加课程", func(t *testing.T) {
		// 添加课程的测试代码
	})

	t.Run("查询课程", func(t *testing.T) {
		_, err := ur.CourseByWeekUsername("1234567", pub.HRBUST, 1)
		assert.Nil(t, err, "查询课程失败")
	})

	t.Run("查询成绩", func(t *testing.T) {
		gpa1, gpa2 := ur.WeightedAverage("2204010417", pub.HRBUST, pub.UG)
		assert.Nil(t, err, "查询成绩失败")
		t.Log(gpa1, gpa2)
	})

	t.Run("添加时间表", func(t *testing.T) {
		times := make([]repository.TimeTable, 0)
		for i, v := range courses {
			times = append(times, repository.TimeTable{
				School:    "hlju",
				Sort:      uint(i + 1),
				StartTime: v.startTime,
				EndTime:   v.endTime,
			})
		}
		assert.Nil(t, ur.AddTimeTable(&times), "添加时间表失败")
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
	{"21:00", "21:45"},
}
