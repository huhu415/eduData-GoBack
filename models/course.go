package models

// Course 为什么要设置周, 而不是利用>BeginWeek && <EndWeek 因为可以方便处理单双周问题
type Course struct {
	ID                    uint   `gorm:"primarykey"`
	StuID                 string `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 学号
	School                string `gorm:"index:idx_stuid_school_week_stype; not null; default:'hrbust'"` // 学校
	Week                  int    `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 周几 没有的话就设置为0
	StuType               int    `gorm:"index:idx_stuid_school_week_stype; not null"`                   // 本科生还是研究生
	WeekDay               int    `gorm:"not null"`                                                      // 星期几 没有的话就设置为0
	NumberOfLessons       int    `gorm:"not null"`                                                      // 第几节课
	NumberOfLessonsLength int    `gorm:"not null"`                                                      // 课程长度
	CourseContent         string `gorm:"not null"`                                                      // 课程名称或内容
	Color                 string `gorm:"not null; default:'#c1d1e0'"`                                   // 课程颜色
	CourseLocation        string
	TeacherName           string
	BeginWeek             int
	EndWeek               int
}

// AddCourse 添加多条课程, 并且把用户名也添加进去
func AddCourse(courses []Course, username, school string, studentType int) {
	// 如果有StuType, 那么添加到结构体中, 不然就不添加
	for index := range courses {
		courses[index].StuID = username
		courses[index].School = school
		courses[index].StuType = studentType
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
	db.Where("stu_id = ? AND school = ? AND week = ?", username, school, week).Find(&courses)
	return courses
}
