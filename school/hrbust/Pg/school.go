package hrbustPg

import (
	"eduData/repository"
	"eduData/school"
	"eduData/school/pub"
	"errors"
	"net/http/cookiejar"
)

const LEFTCORUSEALL = "Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID="

type HrbustPg struct {
	stuID  string
	passWd string
	cookie *cookiejar.Jar
}

func NewHrbustPg(stuID, passWd string, c ...*cookiejar.Jar) school.School {
	h := HrbustPg{
		stuID:  stuID,
		passWd: passWd,
	}
	if len(c) == 1 {
		h.cookie = c[0]
	}
	return &h
}

func (h *HrbustPg) SetCookie(c *cookiejar.Jar) {
	h.cookie = c
}

func (h *HrbustPg) Cookie() *cookiejar.Jar {
	return h.cookie
}
func (h *HrbustPg) SchoolName() pub.SchoolName {
	return pub.HRBUST
}
func (h *HrbustPg) StuType() pub.StuType {
	return pub.PG
}
func (h *HrbustPg) StuID() string {
	return h.stuID
}
func (h *HrbustPg) PassWd() string {
	return h.passWd
}

func (h *HrbustPg) Signin() error {
	c, err := Signin(h.stuID, h.passWd)
	if err != nil {
		return err
	}
	h.cookie = c
	return nil
}

func (h *HrbustPg) GetCourse() ([]repository.Course, error) {
	pgHTML, errPg := GetData(h.cookie, h.stuID, LEFTCORUSEALL)
	if errPg != nil {
		return nil, errPg
	}

	return ParseDataCouresAll(pgHTML)
}

func (h *HrbustPg) GetGrade() ([]repository.CourseGrades, error) {
	return nil, errors.New("not implement")
}
