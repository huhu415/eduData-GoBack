package models

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

func AddCourseGrades(CourseGrades []CourseGrades, username string) {
	for index := range CourseGrades {
		CourseGrades[index].StuID = username
	}
	db.Create(&CourseGrades)
}

func DeleteUserAllCourseGrades(username, school string) {
	db.Where("stu_id = ? AND school = ?", username, school).Delete(&CourseGrades{})
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
func WeightedAverage(username, school string, stuType int) (float64, float64) {
	var result1, result2 float64
	db.Raw("SELECT round(COALESCE(SUM( course_grade * course_credit ), 0) / COALESCE(SUM ( course_credit ), 1),2) FROM course_grades WHERE course_type = '必修' AND course_grade >= 60  AND course_credit != 0 AND stu_id = ? AND school = ? and stu_type = ?", username, school, stuType).Scan(&result1)
	db.Raw("SELECT round(COALESCE(sum((CASE WHEN course_grade >= 80 THEN 4.0 WHEN course_grade >= 70 THEN 3.0 WHEN course_grade >= 60 THEN 2.0 WHEN course_grade >= 50 THEN 1.0 ELSE 0.0 END)* course_credit), 0)/ COALESCE(sum(course_credit), 1),2)FROM course_grades WHERE course_credit != 0 AND course_grade!= 0 AND stu_id = ? AND school = ? and stu_type = ?", username, school, stuType).Scan(&result2)
	return result1, result2
}
