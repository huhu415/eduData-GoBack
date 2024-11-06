package hljuUg

import (
	"net/http/cookiejar"
	"time"

	"eduData/repository"
	"eduData/school"
	"eduData/school/pub"

	"github.com/sirupsen/logrus"
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
	xn, xq := h.getYearTerm()
	d, err := GetData(h.cookie, xn, xq)
	if err != nil {
		return nil, err
	}

	return ParseCoruse(d)
}

func (h *HljuUg) GetGrade() ([]repository.CourseGrades, error) {
	s, err := GetScore(h.cookie)
	if err != nil {
		return nil, err
	}
	return ParseScore(s)
}

func (h *HljuUg) getYearTerm() (int, int) {
	now := time.Now()
	nowYear, nowMonth := now.Year(), now.Month()
	xn, xq := 0, 0

	// 判断学期和学年
	if nowMonth >= 9 { // 9 月到 12 月为第一学期
		xn, xq = nowYear, 1
	} else if nowMonth >= 3 { // 3 月到 8 月为第二学期
		xn, xq = nowYear-1, 2
	} else { // 1 月到 2 月仍在上一学年的第一学期
		xn, xq = nowYear-1, 1
	}

	logrus.Debugf("year: %d, term: %d", xn, xq)
	return xn, xq
}
