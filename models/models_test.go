package models

import (
	"testing"
)

// TestNewDatabase 测试NewDatabase是否能够连接上数据库
func TestNewDatabase(t *testing.T) {
	NewDatabase()
}

// TestAddCourse 测试AddCourse是否能够添加课程
func TestAddCourse(t *testing.T) {
	NewDatabase()
	//AddCourse()
}

// TestCourseByWeekUsername 测试CourseByWeekUsername是否能够查询课程
func TestCourseByWeekUsername(t *testing.T) {
	NewDatabase()
	courses := CourseByWeekUsername(1, "1234567", "hrbust")
	t.Log(courses)
}

// 测试能否计算出来绩点
func TestCalculateGPA(t *testing.T) {
	NewDatabase()
	gpa1, gpa2 := WeightedAverage("2204010417", "hrbust", 1)
	t.Log(gpa1, gpa2)
}
