package models

import (
	"github.com/sirupsen/logrus"
)

type StuInfo struct {
	ID      uint   `gorm:"primarykey"`                      // 主键
	StuID   string `gorm:"uniqueIndex:id_school; not null"` // 学号
	School  string `gorm:"uniqueIndex:id_school; not null"` // 学校
	StuType int    `gorm:"uniqueIndex:id_school; not null"` // 本科生还是研究生

	Name      string // 姓名
	Gender    string // 性别
	Grade     string // 年级
	Phone     string // 联系电话
	IDCard    string // 身份证号
	College   string // 学院
	Major     string // 专业
	Class     string // 班级
	Apartment string // 所在公寓
	RoomFloor string // 所在楼层
	RoomNum   string // 所在房间号
}

// 如果有记录则更新，没有则插入
func (s *StuInfo) CreatAndUpdataStuInfo() error {
	var stu StuInfo
	db.Where("stu_id = ? AND school = ? AND stu_type = ?", s.StuID, s.School, s.StuType).
		Find(&stu).Limit(1)
	if stu.ID == 0 {
		logrus.Infof("\033[1;32m 插入学生信息%+v \033[0m", s)
		return db.Create(s).Error
	}
	logrus.Infof("\033[1;32m 更新学生信息%+v \033[0m", stu)
	return db.Model(&stu).Updates(s).Error
}
