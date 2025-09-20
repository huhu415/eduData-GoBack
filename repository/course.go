package repository

import (
	"eduData/school/pub"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Course 为什么要设置周, 而不是利用>BeginWeek && <EndWeek 因为可以方便处理单双周问题
type Course struct {
	ID                    uint           `gorm:"primarykey"`
	StuID                 string         `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 学号
	School                pub.SchoolName `gorm:"index:idx_stuid_school_week_stype; not null; default:'hrbust'"` // 学校
	Week                  int            `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 第几周 没有的话就设置为0
	Weeks                 pq.Int64Array  `gorm:"column:weeks;type:integer[]"`                                   // 一个课程哪周有
	StuType               pub.StuType    `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 本科生还是研究生
	WeekDay               int            `gorm:"not null"`                                                      // 星期几 没有的话就设置为0
	NumberOfLessons       int            `gorm:"not null"`                                                      // 第几节课
	NumberOfLessonsLength int            `gorm:"not null"`                                                      // 课程长度
	CourseContent         string         `gorm:"not null"`                                                      // 课程名称或内容
	Color                 string         `gorm:"not null; default:'#c1d1e0'"`                                   // 课程颜色
	CourseLocation        string
	CourseLocations       pq.StringArray `gorm:"column:course_locations;type:text[]"` // 课程地点
	TeacherName           string
	BeginWeek             int
	EndWeek               int
}

// AddCourse 添加多条课程, 并且把用户名,学校,研究生还是本科也添加进去
func (r *Repository) AddCourse(Course []Course) error {
	return r.database.Create(&Course).Error
}

// CourseByWeekUsername 通过周数和用户名查询符合的课程
func (r *Repository) CourseByWeekUsername(username string, school pub.SchoolName, week int) ([]Course, error) {
	var course []Course
	err := r.database.
		Where("stu_id = ? AND school = ? AND week = ?", username, school, week).
		Find(&course).Error
	return course, err
}

// 利用事务删除用户的所有课程, 并且添加新的课程
func (r *Repository) DeleteAndCreateCourse(course []Course) error {
	if len(course) == 0 {
		return nil
	}
	return r.database.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("stu_id = ? AND school = ?", course[0].StuID, course[0].School).
			Delete(&Course{}).Error; err != nil {
			return err
		}
		return tx.Create(&course).Error
	})
}

func (r *Repository) CourseByGroup(username string, school pub.SchoolName, stuType pub.StuType) ([]Course, error) {
	/*
			SELECT  week_day,  number_of_lessons, number_of_lessons_length, course_content, course_location, teacher_name, begin_week, end_week, school,count(*)
		from courses
		where stu_id = 'A05250061'
		group by  week_day,  number_of_lessons, number_of_lessons_length, course_content, course_location, teacher_name, begin_week, end_week, school

	*/

	allCourse := "stu_id, school, stu_type, course_content, teacher_name"

	selectCols := allCourse + ", array_agg(week ORDER BY week) AS weeks, array_agg(distinct course_location ORDER BY course_location) AS course_locations"

	var course []Course
	err := r.database.
		Model(&Course{}).
		Where("stu_id = ? AND school = ? AND stu_type = ?", username, school, stuType).
		Select(selectCols).
		Group(allCourse).
		Find(&course).Error
	return course, err
}

func (r *Repository) CourseByCourseTeacher(username string, school pub.SchoolName, stuType pub.StuType, courseName string, teacher string) ([]Course, error) {
	r.database.Debug()

	allCourse := "stu_id, school, stu_type, course_content, teacher_name, course_location, week_day, begin_week, end_week"
	selectCols := allCourse + ", array_agg(week ORDER BY week) AS weeks"

	var course []Course
	err := r.database.
		Model(&Course{}).
		Where("stu_id = ? AND school = ? AND stu_type = ? AND course_content = ? AND teacher_name = ?",
			username, school, stuType, courseName, teacher).
		Select(selectCols).
		Group(allCourse).
		Find(&course).Error
	return course, err
}
