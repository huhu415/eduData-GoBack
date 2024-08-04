// Package app 处理路由的逻辑
package app

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"

	"eduData/api/middleware"
	"eduData/api/pub"
	"eduData/domain"
	"eduData/models"
)

// Signin 下发jwt的cookie
func Signin(c *gin.Context) {
	le, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	//创建jwt
	j := middleware.NewJWT()
	tokenString, err := j.CreateToken(jwt.MapClaims{
		"username": loginForm.Username,
	})
	if err != nil {
		le.Errorf("can not jwt createToken %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("请重试 %w", err)).Error(),
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
	le, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	cookieAny, ok := c.Get("cookie")
	if !ok {
		le.Error("signin中间件没有正常设置cookie")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("请重试 %w", err)).Error(),
		})
		return
	}
	cookie := cookieAny.(*cookiejar.Jar)

	//删除数据库中的所有这个用户名的课程
	models.DeleteUserAllCourse(loginForm.Username, loginForm.School)

	table, err := pub.JudgeUgOrPgGetInfo(loginForm, cookie)
	if err != nil {
		le.Errorf("解析课程出错 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("获取课程错误 %w", err)).Error(),
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
	le, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	cookieAny, ok := c.Get("cookie")
	if !ok {
		le.Error("signin中间件没有正常设置cookie")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("请重试 %w", err)).Error(),
		})
		return
	}
	Cookiejar := cookieAny.(*cookiejar.Jar)

	//删除数据库中的所有这个用户名的课程
	models.DeleteUserAllCourseGrades(loginForm.Username, loginForm.School)

	grade, err := pub.JudgeUgOrPgGetGrade(loginForm, Cookiejar)
	if err != nil {
		le.Errorf("获取成绩错误 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("获取成绩错误 %w", err)).Error(),
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
	le, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	//获取url中的周数的参数
	week := c.Param("week")
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		le.Errorf("周数必须是数字 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("周数必须是数字 %w", err)).Error(),
		})
		return
	} else if !(weekInt >= 0 && weekInt <= 30) {
		le.Errorf("周数必须在0-30之间")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("周数必须在0-30之间 %w", err)).Error(),
		})
		return
	}
	le.Debugf("GetWeekCoure()中的 week:%d\n", weekInt)

	//获取数据库中的课程
	courseByWeekUsername := models.CourseByWeekUsername(weekInt, loginForm.Username, loginForm.School)
	c.JSON(http.StatusOK, courseByWeekUsername)
}

// GetGrade 获取成绩
func GetGrade(c *gin.Context) {
	_, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
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
	_, loginForm, err := pub.GetLogerEntryANDLoginForm(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeTable": models.GetTimeTable(loginForm.School),
	})
}

// AddCoures 增加课程
func AddCoures(c *gin.Context) {
	var data domain.AddcouresStruct
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": c.Error(fmt.Errorf("格式错误, %w", err)).Error(),
		})
		return
	}

	models.AddCourse(pub.ParseAddCrouse(&data), data.LoginForm.Username, data.LoginForm.School, data.LoginForm.StudentType)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程添加成功",
	})
}
