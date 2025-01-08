package hljuUg

import (
	"encoding/json"
	"strconv"
	"strings"

	"eduData/repository"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Code    int     `json:"code"`    // 状态码
	Msg     *string `json:"msg"`     // 消息，可能为空
	MsgEn   *string `json:"msg_en"`  // 英文消息，可能为空
	Content Content `json:"content"` // 内容部分
}

// Content 代表响应中的内容部分
// 包含课程列表及分页信息
type Content struct {
	Total             int      `json:"total"`             // 课程总数
	List              []Course `json:"list"`              // 课程列表
	PageNum           int      `json:"pageNum"`           // 当前页码
	PageSize          int      `json:"pageSize"`          // 每页大小
	Size              int      `json:"size"`              // 当前页的实际大小
	StartRow          int      `json:"startRow"`          // 当前页的开始行号
	EndRow            int      `json:"endRow"`            // 当前页的结束行号
	Pages             int      `json:"pages"`             // 总页数
	PrePage           int      `json:"prePage"`           // 上一页页码
	NextPage          int      `json:"nextPage"`          // 下一页页码
	IsFirstPage       bool     `json:"isFirstPage"`       // 是否为第一页
	IsLastPage        bool     `json:"isLastPage"`        // 是否为最后一页
	HasPreviousPage   bool     `json:"hasPreviousPage"`   // 是否有上一页
	HasNextPage       bool     `json:"hasNextPage"`       // 是否有下一页
	NavigatePages     int      `json:"navigatePages"`     // 导航页数
	NavigatepageNums  []int    `json:"navigatepageNums"`  // 导航页码列表
	NavigateFirstPage int      `json:"navigateFirstPage"` // 导航的第一页页码
	NavigateLastPage  int      `json:"navigateLastPage"`  // 导航的最后一页页码
	Sfxsfx            string   `json:"sfxsfx"`            // 额外信息
}

// Course 代表课程信息的结构体
// 包含课程的各项属性
type Course struct {
	Xs       string  `json:"xs"`        // 学时
	Kcmc     string  `json:"kcmc"`      // 课程名称
	Kcdm     string  `json:"kcdm"`      // 课程代码
	Xnxq     string  `json:"xnxq"`      // 学年学期
	Xnxqmc   string  `json:"xnxqmc"`    // 学年学期名称
	Xnxqmcen string  `json:"xnxqmcen"`  // 学年学期英文名称
	KcmcEn   *string `json:"kcmc_en"`   // 课程英文名称，可能为空
	Kcxz     string  `json:"kcxz"`      // 课程性质
	Kcxzen   *string `json:"kcxzen"`    // 课程性质英文名称，可能为空
	Kclb     string  `json:"kclb"`      // 课程类别
	Kclben   *string `json:"kclben"`    // 课程类别英文名称，可能为空
	Zzcj     string  `json:"zzcj"`      // 最终成绩
	Zzcjen   *string `json:"zzcjen"`    // 最终成绩英文，可能为空
	Yxmc     string  `json:"yxmc"`      // 院系名称
	Zzzscj   string  `json:"zzzscj"`    // 最终折算成绩
	Bkcx     string  `json:"bkcx"`      // 补考情况
	BkcxEn   *string `json:"bkcx_en"`   // 补考情况英文，可能为空
	Khfs     string  `json:"khfs"`      // 考核方式
	KhfsEn   *string `json:"khfs_en"`   // 考核方式英文，可能为空
	Cjbzmc   *string `json:"cjbzmc"`    // 成绩备注，可能为空
	Xf       float64 `json:"xf"`        // 学分
	Xscj     string  `json:"xscj"`      // 学生成绩
	Xscjen   *string `json:"xscjen"`    // 学生成绩英文，可能为空
	Xszscj   string  `json:"xszscj"`    // 学生折算成绩
	Zpcj     string  `json:"zpcj"`      // 总评成绩
	Zpzscj   string  `json:"zpzscj"`    // 总评折算成绩
	Sfkcx    *string `json:"sfkcx"`     // 是否考查，可能为空
	Sfjg     *string `json:"sfjg"`      // 是否合格，可能为空
	Sfyfx    string  `json:"sfyfx"`     // 是否有效
	Id       string  `json:"id"`        // 课程ID
	Rwid     string  `json:"rwid"`      // 任务ID
	Glcjid   string  `json:"glcjid"`    // 关联成绩ID
	Rwh      string  `json:"rwh"`       // 任务号
	Sfpjcxcj string  `json:"sfpjcxcj"`  // 是否评价下次成绩
	Sfxsxq   *string `json:"sfxsxq"`    // 是否显示学期，可能为空
	Pm       string  `json:"pm"`        // 排名
	Zrs      string  `json:"zrs"`       // 总人数
	Xnxqx    *string `json:"xnxqx"`     // 学年学期，可能为空
	Sfmx     *string `json:"sfmx"`      // 是否明细，可能为空
	YxmcEn   *string `json:"yxmc_en"`   // 院系英文名称，可能为空
	CjbzmcEn *string `json:"cjbzmc_en"` // 成绩备注英文，可能为空
	Xscjlb   *string `json:"xscjlb"`    // 成绩类别，可能为空
	Sffx     *string `json:"sffx"`      // 是否复选，可能为空
	Shzt     *string `json:"shzt"`      // 审核状态，可能为空
}

func ParseScore(b *[]byte) ([]repository.CourseGrades, error) {
	var resGrades []repository.CourseGrades

	var response Response
	err := json.Unmarshal(*b, &response)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	for _, course := range response.Content.List {
		var gradeNum float64
		if course.Zzcj != "" {
			gradeNum, err = strconv.ParseFloat(course.Zzcj, 64)
			if err != nil {
				// 解析不了的话就是-1
				gradeNum = -1
			}
		}
		lenXnxq := len(course.Xnxq)
		// 2020-20211 我要分离出最后一位
		yearNum, semesterNum := course.Xnxq[:lenXnxq-1], course.Xnxq[lenXnxq-1]
		courseType := strings.TrimSuffix(course.Kcxz, "课")

		grade := repository.CourseGrades{
			// StuID        string         `gorm:"index; not null"` // 学号
			// School       pub.SchoolName `gorm:"index; not null"` // 学校
			// StuType      pub.StuType    `gorm:"not null"`        // 本科生还是研究生
			Year:         yearNum,
			Semester:     string(semesterNum),
			CourseName:   course.Kcmc,
			CourseType:   courseType,
			CourseCredit: float64(course.Xf),
			CourseGrade:  gradeNum,
		}
		resGrades = append(resGrades, grade)
	}

	return resGrades, nil
}
