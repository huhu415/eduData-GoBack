package models

// TimeTable 时间表, 可以实现不同学校, 不同年级, 不同的时间表
type TimeTable struct {
	ID        uint   `gorm:"primarykey"`       // 主键
	School    string `gorm:"index; not null;"` // 学校
	Sort      uint   `gorm:"not null;"`        // 排序
	StartTime string `gorm:"not null;"`        // 开始时间
	EndTime   string `gorm:"not null;"`        // 结束时间
	grade     string // 年级
}

// GetTimeTable 通过学校获取时间表
func GetTimeTable(school string) []TimeTable {
	var timeTables []TimeTable
	db.Where("school = ?", school).Order("sort").Find(&timeTables)
	return timeTables
}
