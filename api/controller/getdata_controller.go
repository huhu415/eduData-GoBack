package controller

import (
	"eduData/domain"
	"eduData/pub"
	"eduData/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDataController struct {
	LoginUsecase *usecase.Usecase
}

// GetWeekCoure 获取某周课程表
func (g *GetDataController) GetWeekCoure(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	courseByWeekUsername, err := g.LoginUsecase.GetCourseByWeek(s, c.Param("week"))
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
func (g *GetDataController) GetGrade(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	// 获取数据库中所有成绩, 还有加权平均分和绩点
	CourseGradesByUsername, CourseGradesPrompt, WeightedAverage, AcademicCredits := g.LoginUsecase.GetGrade(s)

	// 返回给前端
	c.JSON(http.StatusOK, gin.H{
		"WeightedAverage":    WeightedAverage,
		"AcademicCredits":    AcademicCredits,
		"CourseGradesPrompt": CourseGradesPrompt,
		"CourseGrades":       CourseGradesByUsername,
	})
}

func (g *GetDataController) GetTimeTable(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	res, err := g.LoginUsecase.GetTimeTable(s)
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

func (g *GetDataController) AddCoures(c *gin.Context) {
	s, _, err := pub.GetSchoolAndLogrus(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.Response{
			Status: domain.FAIL,
			Msg:    c.Error(err).Error(),
		})
		return
	}

	var data domain.AddcouresStruct
	if err = g.LoginUsecase.AddCourse(pub.ParseAddCrouse(&data), s); err != nil {
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
