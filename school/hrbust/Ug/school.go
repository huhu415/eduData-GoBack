package hrbustUg

import (
	"context"
	"eduData/models"
	school "eduData/school"
	"errors"
	"net/http/cookiejar"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	YEAR = "44"
	TERM = "2"
)

type HrbustUg struct {
	stuID  string
	passWd string
	cookie *cookiejar.Jar
}

func NewHrbustUg(stuID, passWd string, c ...*cookiejar.Jar) school.School {
	h := HrbustUg{
		stuID:  stuID,
		passWd: passWd,
	}
	if len(c) == 1 {
		h.cookie = c[0]
	}
	return &h
}

func (h *HrbustUg) SchoolName() string {
	return "hrbust"
}
func (h *HrbustUg) StuType() int {
	return 1
}
func (h *HrbustUg) StuID() string {
	return h.stuID
}
func (h *HrbustUg) PassWd() string {
	return h.passWd
}
func (h *HrbustUg) Cookie() *cookiejar.Jar {
	return h.cookie
}

func (h *HrbustUg) Signin() error {
	c, err := Signin(h.stuID, h.passWd)
	if err != nil {
		return err
	}
	h.cookie = c
	return nil
}

func (h *HrbustUg) GetCourse() ([]models.Course, error) {
	b, err := GetCourseByTime(h.cookie, YEAR, TERM)
	if err != nil {
		return nil, err
	}
	return ParseDataCrouseAll(b)
}

// YearSemester 年与学期的结构体
type yearSemester struct {
	Year     string // 43是23年, 44是24年
	Semester string // 1是春季-下学期, 2是秋季-上学期
}

func (h *HrbustUg) GetGrade() ([]models.CourseGrades, error) {
	if h.cookie == nil {
		return nil, errors.New("not found the cookie")
	}
	// 3个协程获取成绩
	var grade []models.CourseGrades
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errs, ctx := errgroup.WithContext(ctx)
	msg := make(chan yearSemester, 10)
	var mutex sync.Mutex
	for i := 0; i < 3; i++ {
		errs.Go(func() error {
			for data := range msg {
				// 获取页面
				ugHTML, errUg := GetDataScore(h.cookie, data.Year, data.Semester)
				if errUg != nil {
					return errUg
				}

				//解析页面, 获得成绩
				table, errUg := ParseDataSore(ugHTML, data.Year, data.Semester)
				if errUg != nil {
					return errUg
				}

				mutex.Lock()
				grade = append(grade, table...)
				mutex.Unlock()
			}
			return nil
		})
	}
	// 添加任务
	atoiYear, err := strconv.Atoi("20" + h.stuID[0:2])
	if err != nil {
		return nil, err
	}
	for i := atoiYear; i <= time.Now().Year(); i++ {
		if i != atoiYear {
			// 第一年没有春季成绩, 所以不是第一年的时候才添加春季
			msg <- yearSemester{Year: strconv.Itoa(i%100 + 20), Semester: "1"}
		}
		msg <- yearSemester{Year: strconv.Itoa(i%100 + 20), Semester: "2"}
	}

	close(msg)
	return grade, errs.Wait()
}

func (h *HrbustUg) GetTimetable() ([]models.TimeTable, error) {
	return models.GetTimeTable(h.SchoolName()), nil
}