// Package app 处理路由的逻辑
package app

import (
	hrbustUg "eduData/school/hrbust/Ug"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"eduData/api/middleware"
	"eduData/domain"
	"eduData/models"
	"eduData/pub"
)

// Signin 下发jwt的cookie
func Signin(c *gin.Context) {
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	if s.SchoolName() == "hrbust" && s.StuType() == 1 {
		body, err := hrbustUg.GetUserInfo(s.Cookie())
		if err != nil {
			le.Errorf("获取个人信息失败 %v", err)
		}

		stuInfo, err := hrbustUg.ParseDataPersonalInfo(body)
		if err != nil {
			le.Errorf("解析个人信息出错 %v", err)
		}

		_ = stuInfo.CreatAndUpdataStuInfo()
	}

	//创建jwt
	j := middleware.NewJWT()
	tokenString, err := j.CreateToken(jwt.MapClaims{
		"school":   s.SchoolName(),
		"username": s.StuID(),
		"stutype":  s.StuType(),
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
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	//删除数据库中的所有这个用户名的课程
	models.DeleteUserAllCourse(s.StuID(), s.SchoolName())

	// 获取课程
	course, err := s.GetCourse()
	if err != nil {
		le.Errorf("获取课程错误 %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	//把解析到的课程与学号都一起添加到数据库中
	models.AddCourse(course, s.StuID(), s.SchoolName(), s.StuType())

	// 成功返回
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程更新成功",
	})
}

// UpdataGrade 更新数据库中的成绩
func UpdataGrade(c *gin.Context) {
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	//删除数据库中的所有这个用户名的课程
	models.DeleteUserAllCourseGrades(s.StuID(), s.SchoolName())

	// 获取成绩
	grade, err := s.GetGrade()
	if err != nil {
		le.Errorf("获取成绩错误 %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	models.AddCourseGrades(grade, s.StuID())

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "成绩更新成功",
	})
}

// GetWeekCoure 获取某周课程表
func GetWeekCoure(c *gin.Context) {
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
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
	courseByWeekUsername := models.CourseByWeekUsername(weekInt, s.StuID(), s.SchoolName())
	c.JSON(http.StatusOK, courseByWeekUsername)
}

// GetGrade 获取成绩
func GetGrade(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	// 获取数据库中所有成绩, 和有那些组
	CourseGradesByUsername, CourseGradesPrompt := models.CourseGradesByUsername(s.StuID(), s.SchoolName())
	// 计算加权平均分
	WeightedAverage, AcademicCredits := models.WeightedAverage(s.StuID(), s.SchoolName(), s.StuType())

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
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"timeTable": models.GetTimeTable(s.SchoolName()),
	})
}

// AddCoures 增加课程
func AddCoures(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": c.Error(err).Error(),
		})
		return
	}

	var data domain.AddcouresStruct
	models.AddCourse(pub.ParseAddCrouse(&data), s.StuID(), s.SchoolName(), s.StuType())

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "课程添加成功",
	})
}
