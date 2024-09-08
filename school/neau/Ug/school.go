package neauUg

import (
	"eduData/repository"
	"eduData/school"
	"eduData/school/pub"
	"errors"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func (n *NeauUg) getYearTerm() (yearTerm string) {
	now := time.Now()
	nowYear, nowMonth := now.Year(), now.Month()
	var resYear, resTerm int
	if nowMonth >= 2 && nowMonth <= 7 {
		resTerm, resYear = 2, nowYear-1
	} else {
		resTerm, resYear = 1, nowYear
	}

	yearTerm = strconv.Itoa(resYear) + "-" + strconv.Itoa(resYear+1) + "-" + strconv.Itoa(resTerm)
	logrus.Debugf("yearTerm: %s", yearTerm)
	return
}

type NeauUg struct {
	stuID  string
	passWd string
	cookie *cookiejar.Jar
}

func NewNeauUg(stuID, passWd string, c ...*cookiejar.Jar) school.School {
	n := NeauUg{
		stuID:  stuID,
		passWd: passWd,
	}
	if len(c) == 1 {
		n.cookie = c[0]
	}
	return &n
}

func (n *NeauUg) SetCookie(c *cookiejar.Jar) {
	n.cookie = c
}

func (n *NeauUg) SchoolName() pub.SchoolName {
	return pub.NEAU
}
func (n *NeauUg) StuType() pub.StuType {
	return pub.UG
}
func (n *NeauUg) StuID() string {
	return n.stuID
}
func (n *NeauUg) PassWd() string {
	return n.passWd
}
func (n *NeauUg) Cookie() *cookiejar.Jar {
	return n.cookie
}

func (n *NeauUg) Signin() error {
	cookie, err := Signin(n.stuID, n.passWd)
	if err != nil {
		return err
	}
	n.cookie = cookie
	return nil
}

func (n *NeauUg) GetCourse() ([]repository.Course, error) {
	GetJSONneau, errNeau := GetData(n.cookie, n.getYearTerm())
	if errNeau != nil {
		return nil, errNeau
	}
	
	return ParseData(GetJSONneau)
}

func (n *NeauUg) GetGrade() ([]repository.CourseGrades, error) {
	return nil, errors.New("不支持这个学校获取成绩")
}
