// Package app 处理路由的逻辑
package app

import (
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"

	"eduData/api/middleware"
	"eduData/domain"
	"eduData/models"
)

// Signin 下发jwt的cookie
func Signin(c *gin.Context) {
	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.Signin()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		return
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
	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.UpdataDB()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		return
	}

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

	//删除数据库中的所有这个用户名的课程
	models.DeleteUserAllCourse(loginForm.Username, loginForm.School)

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
	models.AddCourse(table, loginForm.Username, loginForm.School, loginForm.StudentType)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程更新成功",
	})
}

// UpdataGrade 更新数据库中的成绩
func UpdataGrade(c *gin.Context) {
	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.UpdataGrade()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		return
	}

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

	//删除数据库中的所有这个用户名的课程
	models.DeleteUserAllCourseGrades(loginForm.Username, loginForm.School)

	grade, err := judgeUgOrPgGetGrade(loginForm, Cookiejar)
	if err != nil {
		_ = c.Error(errors.New("app.UpdataGrade()函数中judgeUgOrPgGetGrade的错误: " + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "获取成绩错误" + err.Error(),
		})
		return
	}

	models.AddCourseGrades(grade, loginForm.Username)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "成绩更新成功",
	})
}

// GetWeekCoure 获取某周课程表
func GetWeekCoure(c *gin.Context) {
	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.UpdataGrade()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		return
	}

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
	//fmt.Printf("GetWeekCoure()中的username: %s, week:%d\n", username, weekInt)

	//获取数据库中的课程
	courseByWeekUsername := models.CourseByWeekUsername(weekInt, loginForm.Username, loginForm.School)
	c.JSON(http.StatusOK, courseByWeekUsername)
}

// GetGrade 获取成绩
func GetGrade(c *gin.Context) {
	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.GetGrade()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		return
	}

	// 获取数据库中所有成绩, 和有那些组
	CourseGradesByUsername, CourseGradesPrompt := models.CourseGradesByUsername(loginForm.Username, loginForm.School)
	// 计算加权平均分
	WeightedAverage, AcademicCredits := models.WeightedAverage(loginForm.Username, loginForm.School, loginForm.StudentType)

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
	var loginForm domain.LoginForm
	if err := c.ShouldBindBodyWith(&loginForm, binding.JSON); err != nil {
		_ = c.Error(errors.New("app.GetTimeTable()函数中ShouldBindBodyWith():" + err.Error())).SetType(gin.ErrorTypePrivate)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "表单格式错误,重新登陆后重新提交",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeTable": models.GetTimeTable(loginForm.School),
	})
}

func AddCoures(c *gin.Context) {
	var data domain.AddcouresStruct
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		log.Printf("Error: %#v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "输入格式不符",
		})
		return
	}

	models.AddCourse(parseAddCrouse(&data), data.LoginForm.Username, data.LoginForm.School, data.LoginForm.StudentType)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程添加成功",
	})
}
