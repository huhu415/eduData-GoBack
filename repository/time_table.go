package repository

import (
	"eduData/school/pub"
)

// TimeTable 时间表, 可以实现不同学校, 不同年级, 不同的时间表
type TimeTable struct {
	ID        uint           `gorm:"primarykey"`                      // 主键
	School    pub.SchoolName `gorm:"index:idx_school_sort;not null;"` // 学校
	Sort      uint           `gorm:"index:idx_school_sort;not null;"` // 排序
	StartTime string         `gorm:"not null;"`                       // 开始时间
	EndTime   string         `gorm:"not null;"`                       // 结束时间
	grade     string         // 年级
}

// AddTimeTable 添加时间表
func (r *Repository) AddTimeTable(timeTables *[]TimeTable) error {
	return r.database.Create(timeTables).Error
}

// GetTimeTable 通过学校获取时间表
func (r *Repository) GetTimeTable(school pub.SchoolName) ([]TimeTable, error) {
	var timeTables []TimeTable
	err := r.database.Where("school = ?", school).
		Order("sort").Find(&timeTables).Error
	return timeTables, err
}
