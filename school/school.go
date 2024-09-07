package school

import (
	"eduData/repository"
	"eduData/school/pub"
	"net/http/cookiejar"
)

type School interface {
	SetCookie(*cookiejar.Jar) // 设置cookie

	Cookie() *cookiejar.Jar     // 返回cookie
	SchoolName() pub.SchoolName // 返回学校名
	StuType() pub.StuType       // 返回学生类型
	StuID() string              // 返回学号
	PassWd() string             // 返回密码

	Signin() error                                // 登陆
	GetCourse() ([]repository.Course, error)      // 获取课程
	GetGrade() ([]repository.CourseGrades, error) // 获取成绩

	// GetTimetable() ([]models.TimeTable, error)         // 获取课程时间表
	// GetCourseByWeek(week int) ([]models.Course, error) // 获取某一周课程
	// GetScore() ([]models.CourseGrades, error)          // 获取成绩信息
}
