package pub

import (
	"errors"
	"fmt"

	"eduData/domain"
	"eduData/repository"
	"eduData/school"
	hljuUg "eduData/school/hlju/Ug"
	hrbustUg "eduData/school/hrbust/Ug"
	neauUg "eduData/school/neau/Ug"
	schoolpub "eduData/school/pub"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewSchoolSwitch(d domain.LoginForm) (school.School, error) {
	var s school.School
	switch schoolpub.SchoolName(d.School) {
	// 哈理工
	case schoolpub.HRBUST:
		switch d.StudentType {
		case 1:
			s = hrbustUg.NewHrbustUg(d.Username, d.Password)
		case 2:
			return nil, errors.New(string(d.School) + "研究生登陆功能还未开发")
		}
	// 东北农业大学
	case schoolpub.NEAU:
		switch d.StudentType {
		case 1:
			s = neauUg.NewNeauUg(d.Username, d.Password)
		case 2:
			return nil, errors.New(string(d.School) + "研究生登陆功能还未开发")
		}
	// 黑龙江大学
	case schoolpub.HLJU:
		switch d.StudentType {
		case 1:
			s = hljuUg.NewHljuUg(d.Username, d.Password)
		case 2:
			return nil, errors.New(string(d.School) + "研究生登陆功能还未开发")
		}
	// 其他没有适配的学校
	default:
		return nil, errors.New("不支持的学校")
	}
	return s, nil
}

func GetSchoolAndLogrus(c *gin.Context) (school.School, *logrus.Entry, error) {
	schoolObj, ok := c.Get("SchoolObj")
	if !ok {
		return nil, nil, fmt.Errorf("school obj not found")
	}
	s, ok := schoolObj.(school.School)
	if !ok {
		return nil, nil, fmt.Errorf("school obj convert fail")
	}

	logerEntry, ok := c.Get("logerEntry")
	if !ok {
		return nil, nil, errors.New("logerEntry not found")
	}
	le, ok := logerEntry.(*logrus.Entry)
	if !ok {
		return nil, nil, errors.New("logerEntry convert fail")
	}

	return s, le, nil
}

func ParseAddCrouse(data *domain.AddcouresStruct) []repository.Course {
	var courses []repository.Course
	for _, key := range data.Time {
		course := repository.Course{
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
