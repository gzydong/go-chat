package model

import "time"

type Organize struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 自增ID
	UserId     int       `gorm:"column:user_id;" json:"user_id"`                 // 用户id
	DeptId     int       `gorm:"column:dept_id;" json:"dept_id"`                 // 部门ID
	PositionId int       `gorm:"column:position_id;" json:"position_id"`         // 部门ID
	CreatedAt  time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (Organize) TableName() string {
	return "organize"
}
