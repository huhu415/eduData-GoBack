package models

type StuInfo struct {
	ID      uint   `gorm:"primarykey"`      // 主键
	StuID   string `gorm:"index; not null"` // 学号
	School  string `gorm:"index; not null"` // 学校
	StuType int    `gorm:"not null"`        // 本科生还是研究生

	StuName      string // 姓名
	StuGender    string // 性别
	StuGrade     string // 年级
	StuApartment string // 所在公寓
	StuRoomFloor string // 所在楼层
	StuRoomNum   string // 所在房间号
	StuPhone     string // 联系电话
}
