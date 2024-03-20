package neauUg

import (
	"eduData/School/pub"
	"strings"

	"encoding/json"

	"eduData/database"
)

type TimeAndPlace struct {
	CampusName                   string `json:"campusName"`
	ClassDay                     int    `json:"classDay"`
	ClassSessions                int    `json:"classSessions"`
	ClassWeek                    string `json:"classWeek"`
	ClassroomName                string `json:"classroomName"`
	ContinuingSession            int    `json:"continuingSession"`
	CoureName                    string `json:"coureName"`
	CoureNumber                  string `json:"coureNumber"`
	CoureSequenceNumber          string `json:"coureSequenceNumber"`
	CoursePropertiesName         string `json:"coursePropertiesName"`
	CourseTeacher                string `json:"courseTeacher"`
	ExecutiveEducationPlanNumber string `json:"executiveEducationPlanNumber"`
	ID                           string `json:"id"`
	Kcm                          string `json:"kcm"`
	Sksj                         string `json:"sksj"`
	StudentNumber                string `json:"studentNumber"`
	TeachingBuildingName         string `json:"teachingBuildingName"`
	Time                         string `json:"time"`
	WeekDescription              string `json:"weekDescription"`
	Xf                           string `json:"xf"`
}

type CourseInfo struct {
	AttendClassTeacher   string `json:"attendClassTeacher"`
	CourseCategoryCode   string `json:"courseCategoryCode"`
	CourseCategoryName   string `json:"courseCategoryName"`
	CourseName           string `json:"courseName"`
	CoursePropertiesCode string `json:"coursePropertiesCode"`
	CoursePropertiesName string `json:"coursePropertiesName"`
	DgFlag               string `json:"dgFlag"`
	ExamTypeCode         string `json:"examTypeCode"`
	ExamTypeName         string `json:"examTypeName"`
	Flag                 string `json:"flag"`
	ID                   struct {
		CoureNumber                  string `json:"coureNumber"`
		CoureSequenceNumber          string `json:"coureSequenceNumber"`
		ExecutiveEducationPlanNumber string `json:"executiveEducationPlanNumber"`
		StudentNumber                string `json:"studentNumber"`
	} `json:"id"`
	ProgramPlanName        string         `json:"programPlanName"`
	ProgramPlanNumber      string         `json:"programPlanNumber"`
	RestrictedCondition    string         `json:"restrictedCondition"`
	RlFlag                 string         `json:"rlFlag"`
	SelectCourseStatusCode string         `json:"selectCourseStatusCode"`
	SelectCourseStatusName string         `json:"selectCourseStatusName"`
	Skzcs                  string         `json:"skzcs"`
	StudyModeCode          string         `json:"studyModeCode"`
	StudyModeName          string         `json:"studyModeName"`
	TimeAndPlaceList       []TimeAndPlace `json:"timeAndPlaceList"`
	Unit                   float64        `json:"unit"`
	YwdgFlag               string         `json:"ywdgFlag"`
	Zkxh                   string         `json:"zkxh"`
}
type Schedule struct {
	_        any                     `json:"rwbgLists"`
	AllUnits float64                 `json:"allUnits"`
	Xkxx     []map[string]CourseInfo `json:"xkxx"`
	CSZ      string                  `json:"csz"`
	_        any                     `json:"dateList"`
}

// ParseData 解析本科生课表json格式数据
func ParseData(jsonInfo *[]byte) ([]database.Course, error) {
	// 构造返回参数
	var courses []database.Course

	// 初始化颜色队列
	queue := pub.NewColorList()

	// 解析json
	var schedule Schedule
	err := json.Unmarshal(*jsonInfo, &schedule)
	if err != nil {
		return nil, err
	}

	// 遍历课程信息
	for _, v := range schedule.Xkxx[0] {
		var course database.Course
		course.StuType = 1 // 本科生
		course.School = "neau"
		course.CourseContent = v.CourseName
		replaced := strings.ReplaceAll(v.AttendClassTeacher, "*", "")
		course.TeacherName = strings.ReplaceAll(replaced, " ", "")

		// 如果没有课程的时间或者地点, 则添加一条
		if len(v.TimeAndPlaceList) == 0 {
			// 没有课程就把4要素都置为0
			course.WeekDay = 0
			course.Week = 0
			course.NumberOfLessons = 0
			course.NumberOfLessonsLength = 0

			startWeek, endWeek, _, err := pub.ExtractWeekRange(v.Skzcs)
			if err != nil {
				return nil, err
			}
			course.BeginWeek, course.EndWeek = startWeek, endWeek

			courses = append(courses, course)
		} else {
			course.Color = queue.Remove(queue.Front()).(string)
			for _, timeAndPlace := range v.TimeAndPlaceList {
				course.CourseLocation = timeAndPlace.TeachingBuildingName + timeAndPlace.ClassroomName
				course.WeekDay = timeAndPlace.ClassDay
				course.NumberOfLessons = timeAndPlace.ClassSessions
				course.NumberOfLessonsLength = timeAndPlace.ContinuingSession

				// 匹配单双周和rangWeek
				startWeek, endWeek, evenOrOdd, err := pub.ExtractWeekRange(timeAndPlace.WeekDescription)
				if err != nil {
					return nil, err
				}

				for i := startWeek; i <= endWeek; i++ {
					// 如果是单双周, 则判断是否符合
					if evenOrOdd == 5 {
						course.Week = i
						courses = append(courses, course)
					} else if i%2 == evenOrOdd {
						course.Week = i
						courses = append(courses, course)
					}

				}
			}
		}
	}
	return courses, nil
}
