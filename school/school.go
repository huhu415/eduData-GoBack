package school

import (
	"eduData/models"
	"net/http/cookiejar"
)

type School interface {
	Cookie() *cookiejar.Jar // 返回cookie
	SchoolName() string     // 返回学校名
	StuType() int           // 返回学生类型
	StuID() string          // 返回学号
	PassWd() string         // 返回密码

	Signin() error                            // 登陆
	GetCourse() ([]models.Course, error)      // 获取课程
	GetGrade() ([]models.CourseGrades, error) // 获取成绩

	// GetTimetable() ([]models.TimeTable, error)         // 获取课程时间表
	// GetCourseByWeek(week int) ([]models.Course, error) // 获取某一周课程
	// GetScore() ([]models.CourseGrades, error)          // 获取成绩信息
}
