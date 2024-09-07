package repository

import (
	"eduData/school/pub"
)

// Course 为什么要设置周, 而不是利用>BeginWeek && <EndWeek 因为可以方便处理单双周问题
type Course struct {
	ID                    uint           `gorm:"primarykey"`
	StuID                 string         `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 学号
	School                pub.SchoolName `gorm:"index:idx_stuid_school_week_stype; not null; default:'hrbust'"` // 学校
	Week                  int            `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 周几 没有的话就设置为0
	StuType               pub.StuType    `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 本科生还是研究生
	WeekDay               int            `gorm:"not null"`                                                      // 星期几 没有的话就设置为0
	NumberOfLessons       int            `gorm:"not null"`                                                      // 第几节课
	NumberOfLessonsLength int            `gorm:"not null"`                                                      // 课程长度
	CourseContent         string         `gorm:"not null"`                                                      // 课程名称或内容
	Color                 string         `gorm:"not null; default:'#c1d1e0'"`                                   // 课程颜色
	CourseLocation        string
	TeacherName           string
	BeginWeek             int
	EndWeek               int
}

// AddCourse 添加多条课程, 并且把用户名,学校,研究生还是本科也添加进去
func (r *Repository) AddCourse(Course []Course) error {
	return r.database.Create(&Course).Error
}

// DeleteUserAllCourse 删除用户的所有课程
func (r *Repository) DeleteUserAllCourse(username string, school pub.SchoolName) error {
	return r.database.
		Where("stu_id = ? AND school = ?", username, school).
		Delete(&Course{}).Error
}

// CourseByWeekUsername 通过周数和用户名查询符合的课程
func (r *Repository) CourseByWeekUsername(username string, school pub.SchoolName, week int) ([]Course, error) {
	var course []Course
	err := r.database.
		Where("stu_id = ? AND school = ? AND week = ?", username, school, week).
		Find(&course).Error
	return course, err
}
