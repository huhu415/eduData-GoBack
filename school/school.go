package school

import (
	"eduData/models"
	"net/http/cookiejar"
)

type School interface {
	Signin(userName, passWord string) (*cookiejar.Jar, error)
	GetCourse(cookieJar *cookiejar.Jar, date ...string) ([]models.Course, error)
	GetScore(cookieJar *cookiejar.Jar) ([]models.CourseGrades, error)
	GetTimetable() ([]models.TimeTable, error)
}
