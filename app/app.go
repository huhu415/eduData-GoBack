// Package app 处理路由的逻辑
package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/sync/errgroup"

	hrbustPg "eduData/School/hrbust/Pg"
	hrbustUg "eduData/School/hrbust/Ug"
	neauUg "eduData/School/neau/Ug"
	"eduData/database"
	"eduData/middleware"
)

const (
	// 控制研究生获取页面用的
	// LEFTCORUSE 某一周的
	//LEFTCORUSE = "Course/StuCourseWeekQuery.aspx?EID=vB5Ke2TxFzG4yVM8zgJqaQowdgBb6XLK0loEdeh1pyPrNQM0n6oBLQ==&UID="
	// LEFTCORUSEALL 学期的
	LEFTCORUSEALL = "Course/StuCourseQuery.aspx?EID=pLiWBm!3y8J!emOuKhzHa3uED3OEJzAvyCpKfhbkdg9RKe9VDAjrUw==&UID="
)

// judgeUgOrPgGetInfo 根据学校和研究生本科生判断获取html并解析
func judgeUgOrPgGetInfo(loginForm middleware.LoginForm, cookieJar *cookiejar.Jar) ([]database.Course, error) {
	var table []database.Course
	switch loginForm.School {
	// 哈理工
	case "hrbust":
		switch loginForm.StudentType {
		case 1:
			ugHTML, errUg := hrbustUg.GetData(cookieJar, "2000")
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
			GetJSONneau, errNeau := neauUg.GetData(cookieJar, "2023-2024-2-1") // todo 设计一下获取学期的函数
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

// YearSemester 年与学期的结构体
type YearSemester struct {
	Year     string // 43是23年, 44是24年
	Semester string // 1是春季-下学期, 2是秋季-上学期
}

// judgeUgOrPgGetGrade 根据学校和研究生本科生判断获取成绩的html, 并解析成绩
func judgeUgOrPgGetGrade(loginForm middleware.LoginForm, cookieJar *cookiejar.Jar) ([]database.CourseGrades, error) {
	var grade []database.CourseGrades
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
			msg := make(chan YearSemester, 10)
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
					msg <- YearSemester{Year: strconv.Itoa(i%100 + 20), Semester: "1"}
				}
				msg <- YearSemester{Year: strconv.Itoa(i%100 + 20), Semester: "2"}
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

// Signin 下发jwt的cookie
func Signin(c *gin.Context) {
	var loginForm middleware.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.Signin()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		c.Abort()
	}

	//创建jwt
	j := middleware.NewJWT()
	tokenString, err := j.CreateToken(jwt.MapClaims{
		"username": loginForm.Username,
	})
	if err != nil {
		err = c.Error(errors.New("app.Signin()函数中CreateToken的错误: " + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "请重试",
		})
		return
	}

	//返回给前端结果
	c.SetCookie("authentication", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "登陆成功",
	})
}

// UpdataDB 更新数据库中的课程表
func UpdataDB(c *gin.Context) {
	cookieAny, ok := c.Get("cookie")
	if !ok {
		_ = c.Error(errors.New("app : signin中间件没有正常设置username, cookie, studentType")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "请重试",
		})
		return
	}
	cookie := cookieAny.(*cookiejar.Jar)

	var loginForm middleware.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.UpdataDB()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		c.Abort()
	}

	//删除数据库中的所有这个用户名的课程
	database.DeleteUserAllCourse(loginForm.Username, loginForm.School)

	table, err := judgeUgOrPgGetInfo(loginForm, cookie)
	if err != nil {
		_ = c.Error(errors.New("app.UpdataDB()函数中judgeUgOrPgGetInfo的错误: " + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "获取课表错误" + err.Error(),
		})
		return
	}

	//把解析到的课程与学号都一起添加到数据库中
	database.AddCourse(table, loginForm.Username, loginForm.School, loginForm.StudentType)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程更新成功",
	})

}

// UpdataGrade 更新数据库中的成绩
func UpdataGrade(c *gin.Context) {
	cookieAny, ok := c.Get("cookie")
	if !ok {
		_ = c.Error(errors.New("app : signin中间件没有正常设置username, cookie, studentType")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "请重试",
		})
		return
	}
	Cookiejar := cookieAny.(*cookiejar.Jar)

	var loginForm middleware.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.UpdataGrade()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		c.Abort()
	}

	//删除数据库中的所有这个用户名的课程
	database.DeleteUserAllCourseGrades(loginForm.Username, loginForm.School)

	grade, err := judgeUgOrPgGetGrade(loginForm, Cookiejar)
	if err != nil {
		_ = c.Error(errors.New("app.UpdataGrade()函数中judgeUgOrPgGetGrade的错误: " + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "获取成绩错误" + err.Error(),
		})
		return
	}

	database.AddCourseGrades(grade, loginForm.Username)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "成绩更新成功",
	})

}

// GetWeekCoure 获取某周课程表
func GetWeekCoure(c *gin.Context) {
	//获取url中的周数的参数
	week := c.Param("week")
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		_ = c.Error(errors.New("app.GetWeekCoure():周数必须是数字")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "周数必须是数字",
		})
		return
	} else if !(weekInt >= 0 && weekInt <= 30) {
		_ = c.Error(errors.New("app.GetWeekCoure():周数必须在0-30之间")).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "周数必须在0-30之间",
		})
		return
	}

	var loginForm middleware.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.UpdataGrade()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		c.Abort()
	}
	//fmt.Printf("GetWeekCoure()中的username: %s, week:%d\n", username, weekInt)

	//获取数据库中的课程
	courseByWeekUsername := database.CourseByWeekUsername(weekInt, loginForm.Username, loginForm.School)
	c.JSON(http.StatusOK, courseByWeekUsername)
}

// GetGrade 获取成绩
func GetGrade(c *gin.Context) {
	var loginForm middleware.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.GetGrade()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		c.Abort()
	}

	// 获取数据库中所有成绩, 和有那些组
	CourseGradesByUsername, CourseGradesPrompt := database.CourseGradesByUsername(loginForm.Username, loginForm.School)
	// 计算加权平均分
	WeightedAverage, AcademicCredits := database.WeightedAverage(loginForm.Username, loginForm.School, loginForm.StudentType)

	// 返回给前端
	c.JSON(http.StatusOK, gin.H{
		"WeightedAverage":    WeightedAverage,
		"AcademicCredits":    AcademicCredits,
		"CourseGradesPrompt": CourseGradesPrompt,
		"CourseGrades":       CourseGradesByUsername,
	})
}

// GetTimeTable 获取上课时间
func GetTimeTable(c *gin.Context) {
	var loginForm middleware.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.GetTimeTable()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{
		"timeTable": database.GetTimeTable(loginForm.School),
	})
}

type AddcouresStruct struct {
	middleware.LoginForm
	Color   string      `json:"color"`
	Coures  string      `json:"coures" binding:"required"`
	Teacher string      `json:"teacher"`
	Time    []TimeEntry `json:"time" binding:"required"`
}

type TimeEntry struct {
	Checkboxs  []int  `json:"checkboxs" binding:"required"`
	MultiIndex []int  `json:"multiIndex" binding:"required"`
	Place      string `json:"place"`
}

func AddCoures(c *gin.Context) {
	data := AddcouresStruct{}
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		log.Printf("Error: %#v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "输入格式不符",
		})
		return
	}

	database.AddCourse(ParseAddCrouse(data), data.LoginForm.Username, data.LoginForm.School, data.LoginForm.StudentType)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程添加成功",
	})
}
