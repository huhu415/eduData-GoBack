package hrbustUg

import (
	"eduData/models"
	school "eduData/school"
	"net/http/cookiejar"
)

type HrbustUg struct {
}

func NewHrbustUg() school.School {
	return &HrbustUg{}
}

func (h *HrbustUg) Signin(userName, passWord string) (*cookiejar.Jar, error) {
	return nil, nil
}
func (h *HrbustUg) GetCourse(cookieJar *cookiejar.Jar, date ...string) ([]models.Course, error) {
	return nil, nil
}
func (h *HrbustUg) GetScore(cookieJar *cookiejar.Jar) ([]models.CourseGrades, error) {
	return nil, nil
}
func (h *HrbustUg) GetTimetable() ([]models.TimeTable, error) {
	return nil, nil
}
