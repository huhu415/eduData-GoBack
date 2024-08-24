package pub

import (
	"context"
	"errors"
	"net/http/cookiejar"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	hrbustPg "eduData/school/hrbust/Pg"
	hrbustUg "eduData/school/hrbust/Ug"
	neauUg "eduData/school/neau/Ug"
	"eduData/domain"
	"eduData/models"
)

const (
	// 控制研究生获取页面用的
	// LEFTCORUSE 某一周的
	//LEFTCORUSE = "Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID="
	// LEFTCORUSEALL 学期的
	LEFTCORUSEALL = "Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID="
	NEAUDATA      = "2023-2024-2-1" // todo 需要随学期更改
)

// judgeUgOrPgGetInfo 根据学校和研究生本科生判断获取html并解析
func JudgeUgOrPgGetInfo(loginForm domain.LoginForm, cookieJar *cookiejar.Jar) ([]models.Course, error) {
	var table []models.Course
	switch loginForm.School {
	// 哈理工
	case "hrbust":
		switch loginForm.StudentType {
		case 1:
			// 获取当前学期的课表
			// ugHTML, errUg := hrbustUg.GetData(cookieJar, "2000")
			ugHTML, errUg := hrbustUg.GetCourseByTime(cookieJar, "44", "2")
			if errUg != nil {
				return nil, errUg
			}
			table, errUg = hrbustUg.ParseDataCrouseAll(ugHTML)
			if errUg != nil {
				return nil, errUg
			}
		case 2:
			pgHTML, errPg := hrbustPg.GetData(cookieJar, loginForm.Username, LEFTCORUSEALL)
			if errPg != nil {
				return nil, errPg
			}
			table, errPg = hrbustPg.ParseDataCouresAll(pgHTML)
			if errPg != nil {
				return nil, errPg
			}
		}
	// 东北农业大学
	case "neau":
		switch loginForm.StudentType {
		case 1:
			GetJSONneau, errNeau := neauUg.GetData(cookieJar, NEAUDATA) // todo 设计一下获取学期的函数
			if errNeau != nil {
				return nil, errNeau
			}
			table, errNeau = neauUg.ParseData(GetJSONneau)
			if errNeau != nil {
				return nil, errNeau
			}
		case 2:
			return nil, errors.New(loginForm.School + "研究生登陆功能还未开发")
		}
	// 其他没有适配的学校
	default:
		return nil, errors.New("不支持的学校")
	}
	return table, nil
}

// JudgeSchoolSignIn 判断是哪个学校的用户来登陆
func JudgeSchoolSignIn(loginForm domain.LoginForm) (*cookiejar.Jar, error) {
	switch loginForm.School {
	// 哈理工
	case "hrbust":
		switch loginForm.StudentType {
		case 1:
			return hrbustUg.Signin(loginForm.Username, loginForm.Password)
		case 2:
			return hrbustPg.Signin(loginForm.Username, loginForm.Password)
		}
	// 东北农业大学
	case "neau":
		switch loginForm.StudentType {
		case 1:
			return neauUg.Signin(loginForm.Username, loginForm.Password)
		case 2:
			return nil, errors.New(loginForm.School + "研究生登陆功能还未开发")
		}
	// 其他没有适配的学校
	default:
		return nil, errors.New("不支持的学校")
	}
	return nil, errors.New("不支持的学校")
}

// YearSemester 年与学期的结构体
type yearSemester struct {
	Year     string // 43是23年, 44是24年
	Semester string // 1是春季-下学期, 2是秋季-上学期
}

// JudgeUgOrPgGetGrade 根据学校和研究生本科生判断获取成绩的html, 并解析成绩
// 根据学号的前两位开始, 查询到现在的年份
func JudgeUgOrPgGetGrade(loginForm domain.LoginForm, cookieJar *cookiejar.Jar) ([]models.CourseGrades, error) {
	var grade []models.CourseGrades
	switch loginForm.School {
	// 哈理工
	case "hrbust":
		switch loginForm.StudentType {
		// 本科生
		case 1:
			// 3个协程获取成绩
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			errs, ctx := errgroup.WithContext(ctx)
			msg := make(chan yearSemester, 10)
			var mutex sync.Mutex
			for i := 0; i < 3; i++ {
				errs.Go(func() error {
					for data := range msg {
						// 获取页面
						ugHTML, errUg := hrbustUg.GetDataScore(cookieJar, data.Year, data.Semester)
						if errUg != nil {
							return errUg
						}

						//解析页面, 获得成绩
						table, errUg := hrbustUg.ParseDataSore(ugHTML, data.Year, data.Semester)
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
			atoiYear, err := strconv.Atoi("20" + loginForm.Username[0:2])
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
			if errs.Wait() != nil {
				return nil, errs.Wait()
			}
		case 2:
			return nil, errors.New("不支持研究生")
		default:
			return nil, errors.New("未知学生")
		}
	// 其他没有适配的学校
	default:
		return nil, errors.New("不支持的学校")
	}
	return grade, nil
}

func ParseAddCrouse(data *domain.AddcouresStruct) []models.Course {
	var courses []models.Course
	for _, key := range data.Time {
		course := models.Course{
			Color:                 data.Color,
			TeacherName:           data.Teacher,
			CourseContent:         data.Coures,
			CourseLocation:        key.Place,
			WeekDay:               key.MultiIndex[0],
			NumberOfLessons:       key.MultiIndex[1],
			NumberOfLessonsLength: key.MultiIndex[2],
		}
		// 如果符合read.md中写的情况, 那应该显示先下面
		if course.NumberOfLessons == 0 || course.NumberOfLessonsLength == 0 || course.WeekDay == 0 || key.Checkboxs == nil {
			course.Week = 0
			courses = append(courses, course)
		} else {
			// 哪几周
			for _, keyCheckbos := range key.Checkboxs {
				course.Week = keyCheckbos
				courses = append(courses, course)
			}
		}

	}
	return courses
}

func GetLogerEntryANDLoginForm(c *gin.Context) (*logrus.Entry, domain.LoginForm, error) {
	logerEntry, ok := c.Get("logerEntry")
	if !ok {
		return nil, domain.LoginForm{}, errors.New("服务器内部错误,请联系管理员")
	}
	le, ok := logerEntry.(*logrus.Entry)
	if !ok {
		return nil, domain.LoginForm{}, errors.New("服务器内部错误,请联系管理员")
	}

	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		return nil, domain.LoginForm{}, errors.New("表单格式错误,重新登陆后重新提交")
	}
	return le, loginForm, nil
}
