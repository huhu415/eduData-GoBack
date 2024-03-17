package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"eduData/setting"
)

var db *gorm.DB

// Course 为什么要设置周, 而不是利用>BeginWeek && <EndWeek 因为可以方便处理单双周问题
type Course struct {
	ID                    uint   `gorm:"primarykey"`
	StuID                 string `gorm:"index; not null"`                   // 学号
	School                string `gorm:"index; not null; default:'hrbust'"` // 学校
	StuType               int    `gorm:"not null"`                          // 本科生还是研究生
	Week                  int    `gorm:"not null"`                          // 周几 没有的话就设置为0
	WeekDay               int    `gorm:"not null"`                          // 星期几 没有的话就设置为0
	NumberOfLessons       int    `gorm:"not null"`                          // 第几节课
	NumberOfLessonsLength int    `gorm:"not null"`                          // 课程长度
	CourseContent         string `gorm:"not null"`                          // 课程名称或内容
	CourseLocation        string
	TeacherName           string
	BeginWeek             int
	EndWeek               int
	Color                 string `gorm:"not null; default:'#c1d1e0'"` // 课程颜色
}

// CourseGrades 课程成绩, 还可以计算绩点
type CourseGrades struct {
	ID           uint    `gorm:"primarykey"`      // 主键
	StuID        string  `gorm:"index; not null"` // 学号
	School       string  `gorm:"index; not null"` // 学校
	StuType      int     `gorm:"not null"`        // 本科生还是研究生
	Year         string  `gorm:"not null"`        // 学年
	Semester     string  `gorm:"not null"`        // 学期
	CourseName   string  `gorm:"not null"`        // 课程名称
	CourseType   string  `gorm:"not null"`        // 选修, 任选, 限选, 还是必修
	CourseCredit float64 `gorm:"not null"`        // 学分
	CourseGrade  float64 `gorm:"not null"`        // 成绩
}

func NewDatabase() {
	var err error
	dsn := setting.Postgres
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		println("database : gorm.Open(), database connect error")
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		println("database : db.DB(), database connect error")
		panic(err)
	}

	// Migrate the schema, 创建表用的, 就用一次就完事了
	//if err = db.AutoMigrate(&Course{}, &CourseGrades{}); err != nil {
	//	panic(err)
	//}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(20)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(200)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("database Connected!")
}

// AddCourse 添加多条课程, 并且把用户名也添加进去
func AddCourse(courses []Course, username string) {
	for index := range courses {
		courses[index].StuID = username
	}
	db.Create(&courses)
}

// DeleteUserAllCourse 删除用户的所有课程
func DeleteUserAllCourse(username, school string) {
	db.Where("stu_id = ? AND school = ?", username, school).Delete(&Course{})
}

// CourseByWeekUsername 通过周数和用户名查询符合的课程
func CourseByWeekUsername(week int, username, school string) []Course {
	var courses []Course
	//查询数据
	// select * from courses where week = ? and stu_id = ?
	db.Where("week = ? AND stu_id = ? AND school = ?", week, username, school).Find(&courses)
	return courses
}

func AddCourseGrades(CourseGrades []CourseGrades, username string) {
	for index := range CourseGrades {
		CourseGrades[index].StuID = username
	}
	db.Create(&CourseGrades)
}

func DeleteUserAllCourseGrades(username, school string) {
	db.Where("\"stu_id\" = ? AND \"school\" = ?", username, school).Delete(&CourseGrades{})
}

func CourseGradesByUsername(username, school string) ([]CourseGrades, []CourseGrades) {
	var courseGrades []CourseGrades
	var courseGradesPrompt []CourseGrades
	db.Where("stu_id = ? AND school = ?", username, school).
		Order("year, semester").
		Find(&courseGrades)
	db.Select("year, semester").
		Where("stu_id = ? AND school = ?", username, school).
		Group("year, semester").
		Order("year, semester").
		Find(&courseGradesPrompt)
	return courseGrades, courseGradesPrompt
}

// WeightedAverage 计算加权平均分和加拿大绩点算法
func WeightedAverage(username, school, stuType string) (float64, float64) {
	var result1, result2 float64
	db.Raw("SELECT round(SUM( course_grade * course_credit ) / SUM ( course_credit ),2) FROM course_grades WHERE course_type = '必修' AND course_grade >= 60  AND course_credit != 0 AND stu_id = ? AND school = ? and stu_type = ?", username, school, stuType).Scan(&result1)
	db.Raw("SELECT round(sum((CASE WHEN course_grade >= 80 THEN 4.0 WHEN course_grade >= 70 THEN 3.0 WHEN course_grade >= 60 THEN 2.0 WHEN course_grade >= 50 THEN 1.0 ELSE 0.0 END)* course_credit)/ sum(course_credit),2)FROM course_grades WHERE course_credit != 0 AND course_grade!= 0AND stu_id = ? AND school = ? and stu_type = ?", username, school, stuType).Scan(&result2)
	return result1, result2
}
