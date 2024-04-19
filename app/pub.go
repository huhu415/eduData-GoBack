package app

import (
	"eduData/database"
)

func ParseAddCrouse(data AddcouresStruct) []database.Course {
	var courses []database.Course
	for _, key := range data.Time {
		course := database.Course{
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
