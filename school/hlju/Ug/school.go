package hljuUg

import (
	"eduData/repository"
	"eduData/school"
	"eduData/school/pub"
	"errors"
	"net/http/cookiejar"
)

type HljuUg struct {
	stuID  string
	passWd string
	cookie *cookiejar.Jar
}

func NewHljuUg(stuID, passWd string, c ...*cookiejar.Jar) school.School {
	h := HljuUg{
		stuID:  stuID,
		passWd: passWd,
	}
	if len(c) == 1 {
		h.cookie = c[0]
	}
	return &h
}

func (h *HljuUg) SetCookie(c *cookiejar.Jar) {
	h.cookie = c
}

func (h *HljuUg) SchoolName() pub.SchoolName {
	return pub.HLJU
}

func (h *HljuUg) StuType() pub.StuType {
	// 本科
	return pub.UG
}

func (h *HljuUg) StuID() string {
	return h.stuID
}

func (h *HljuUg) PassWd() string {
	return h.passWd
}

func (h *HljuUg) Cookie() *cookiejar.Jar {
	return h.cookie
}

func (h *HljuUg) Signin() error {
	c, err := Signin(h.stuID, h.passWd)
	if err != nil {
		return err
	}
	h.cookie = c
	return nil
}

func (h *HljuUg) GetCourse() ([]repository.Course, error) {
	d, err := GetData(h.cookie)
	if err != nil {
		return nil, err
	}

	return ParseData(d)
}

func (h *HljuUg) GetGrade() ([]repository.CourseGrades, error) {
	return nil, errors.New("hljuUg.GetGrade() not implement")
}
