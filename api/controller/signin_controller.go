package controller

import (
	"fmt"
	"net/http"
	"time"

	"eduData/api/middleware"
	"eduData/domain"
	"eduData/pub"
	"eduData/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type SigninController struct {
	LoginUsecase *usecase.Usecase
}

func (lc *SigninController) LogIn(c *gin.Context) {
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	if err = lc.LoginUsecase.SigninAndSetCache(s); err != nil {
		le.Errorf("登陆失败 %v", err)
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 创建jwt
	j := middleware.NewJWT()
	tokenString, err := j.CreateToken(jwt.MapClaims{
		"school":   s.SchoolName(),
		"username": s.StuID(),
		"stutype":  s.StuType(),
		"iat":      time.Now().Unix(),
	})
	if err != nil {
		le.Errorf("can not jwt createToken %v", err)
		c.JSON(http.StatusBadRequest, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(fmt.Errorf("请重试 %w", err)).Error(),
		})
		return
	}

	// 返回给前端结果
	c.SetCookie("authentication", tokenString, 3600*24*30, "/", "", false, true)
	c.JSON(http.StatusOK, domain.Response{
		Status: domain.SUCCESS,
		Msg:    "登陆成功",
	})
}

// Updata 更新课程
func (lc *SigninController) UpdateCourse(c *gin.Context) {
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	if err = lc.LoginUsecase.GetAndUpdataCache(s); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 获取课程
	logrus.Debugf("id:%s, schoolName:%s, type:%d", s.StuID(), s.SchoolName(), s.StuType())
	course, err := s.GetCourse()
	if err != nil {
		le.Errorf("获取课程错误 %v", err)
		// 删除缓存
		lc.LoginUsecase.CleanCache(s)
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 把课程添加到数据库, 并且删除原来的课程
	if err = lc.LoginUsecase.DeleteAndCreateCourse(course, s); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 成功返回
	c.JSON(http.StatusOK, domain.Response{
		Status: domain.SUCCESS,
		Msg:    "课程更新成功",
	})
}

// UpdataGrade 更新成绩
func (lc *SigninController) UpdateGrade(c *gin.Context) {
	s, le, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	if err = lc.LoginUsecase.GetAndUpdataCache(s); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	logrus.Debugf("id:%s, schoolName:%s, type:%d", s.StuID(), s.SchoolName(), s.StuType())
	// 获取成绩
	grade, err := s.GetGrade()
	if err != nil {
		le.Errorf("获取成绩错误 %v", err)
		// 删除缓存
		lc.LoginUsecase.CleanCache(s)
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 把成绩添加到数据库, 并且删除原来的成绩
	if err = lc.LoginUsecase.DeleteAndCreateGrade(grade, s); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Status: domain.SUCCESS,
		Msg:    "成绩更新成功",
	})
}

// GetWeekCoure 获取某周课程表
func (lc *SigninController) GetWeekCourse(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	courseByWeekUsername, err := lc.LoginUsecase.GetCourseByWeek(s, c.Param("week"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}
	c.JSON(http.StatusOK, courseByWeekUsername)
}

// GetGrade 获取成绩
func (lc *SigninController) GetGrade(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 获取数据库中所有成绩, 还有加权平均分和绩点
	CourseGradesByUsername, CourseGradesPrompt, WeightedAverage, AcademicCredits := lc.LoginUsecase.GetGrade(s)

	// 返回给前端
	c.JSON(http.StatusOK, gin.H{
		"WeightedAverage":    WeightedAverage,
		"AcademicCredits":    AcademicCredits,
		"CourseGradesPrompt": CourseGradesPrompt,
		"CourseGrades":       CourseGradesByUsername,
	})
}

func (lc *SigninController) GetTimeTable(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	res, err := lc.LoginUsecase.GetTimeTable(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Status: domain.SUCCESS,
		Msg:    res,
	})
}

func (lc *SigninController) AddCourse(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	var data domain.AddcouresStruct
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, domain.Response{
			Status: domain.FAIL,
			Msg:    "无效的请求数据",
		})
		return
	}

	if err = lc.LoginUsecase.AddCourse(pub.ParseAddCrouse(&data), s); err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.Response{
		Status: domain.SUCCESS,
		Msg:    "课程添加成功",
	})
}
