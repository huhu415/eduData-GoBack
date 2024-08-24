package neauUg

import (
	"eduData/models"
	"eduData/school"
	"net/http/cookiejar"
)

type NeauUg struct {
}

func NewNeauUg() school.School {
	return &NeauUg{}
}

func (h *NeauUg) Signin(userName, passWord string) (*cookiejar.Jar, error) {
	return nil, nil
}
func (h *NeauUg) GetCourse(cookieJar *cookiejar.Jar, date ...string) ([]models.Course, error) {
	return nil, nil
}
func (h *NeauUg) GetScore(cookieJar *cookiejar.Jar) ([]models.CourseGrades, error) {
	return nil, nil
}
func (h *NeauUg) GetTimetable() ([]models.TimeTable, error) {
	return nil, nil
}
