package usecase

import (
	"fmt"
	"net/http/cookiejar"
	"strconv"
	"time"

	"eduData/repository"
	"eduData/school"

	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

type Usecase struct {
	repository *repository.Repository
	cache      *cache.Cache
}

func NewUsecase(r *repository.Repository) *Usecase {
	return &Usecase{
		repository: r,
		cache:      cache.New(10*time.Minute, 15*time.Minute),
	}
}

func (u *Usecase) GetAndUpdataCache(s school.School) error {
	stuInfo := fmt.Sprintf("%s-%d-%s", s.SchoolName(), s.StuType(), s.StuID())
	cookie, ok := u.cache.Get(stuInfo)
	if ok {
		logrus.Debugf("从缓存中获取到cookie, syuInfo:%s", stuInfo)
		s.SetCookie(cookie.(*cookiejar.Jar))
	} else {
		if err := s.Signin(); err != nil {
			return err
		}
		u.cache.Set(stuInfo, s.Cookie(), cache.DefaultExpiration)
	}
	return nil
}

func (u *Usecase) CleanCache(s school.School) {
	stuInfo := fmt.Sprintf("%s-%d-%s", s.SchoolName(), s.StuType(), s.StuID())
	u.cache.Delete(stuInfo)
}

func (u *Usecase) SigninAndSetCache(s school.School) error {
	if err := s.Signin(); err != nil {
		return err
	}
	stuInfo := fmt.Sprintf("%s-%d-%s", s.SchoolName(), s.StuType(), s.StuID())
	u.cache.Set(stuInfo, s.Cookie(), cache.DefaultExpiration)
	return nil
}

func (u *Usecase) DeleteAndCreateCourse(Course []repository.Course, school school.School) error {
	for i := range Course {
		Course[i].StuID = school.StuID()
		Course[i].School = school.SchoolName()
		Course[i].StuType = school.StuType()
	}
	return u.repository.DeleteAndCreateCourse(Course)
}

func (u *Usecase) DeleteAndCreateGrade(CourseGrades []repository.CourseGrades, school school.School) error {
	for i := range CourseGrades {
		CourseGrades[i].StuID = school.StuID()
		CourseGrades[i].School = school.SchoolName()
		CourseGrades[i].StuType = school.StuType()
	}
	return u.repository.DeleteAndAddCourseGrades(CourseGrades)
}

func (u *Usecase) GetCourseByWeek(school school.School, week string) ([]repository.Course, error) {
	weekInt, err := strconv.Atoi(week)
	if err != nil {
		return nil, fmt.Errorf("周数必须是数字 %w", err)
	} else if !(weekInt >= 0 && weekInt <= 30) {
		return nil, fmt.Errorf("周数必须在0-30之间")
	}
	return u.repository.CourseByWeekUsername(school.StuID(), school.SchoolName(), weekInt)
}

func (u *Usecase) GetSingleCourseTeacher(school school.School, course string, teacher string) ([]repository.Course, error) {
	return u.repository.CourseByCourseTeacher(school.StuID(), school.SchoolName(), school.StuType(), course, teacher)
}

func (u *Usecase) GetGrade(school school.School) ([]repository.CourseGrades, []repository.CourseGrades, float64, float64) {
	courseGrades, courseGradesPrompt := u.repository.CourseGradesByUsername(school.StuID(), school.StuType(), school.SchoolName())
	WeightedAverage, AcademicCredits := u.repository.WeightedAverage(school.StuID(), school.SchoolName(), school.StuType())
	return courseGrades, courseGradesPrompt, WeightedAverage, AcademicCredits
}

func (u *Usecase) GetTimeTable(school school.School) ([]repository.TimeTable, error) {
	return u.repository.GetTimeTable(school.SchoolName())
}

func (u *Usecase) AddCourse(Course []repository.Course, school school.School) error {
	for index := range Course {
		Course[index].StuID = school.StuID()
		Course[index].School = school.SchoolName()
		Course[index].StuType = school.StuType()
	}
	return u.repository.AddCourse(Course)
}
