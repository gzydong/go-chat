package model

import "time"

type Organize struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 自增ID
	UserId     int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 用户id
	Department string    `gorm:"column:department;default:''" json:"department"`   // 部门ID
	Position   string    `gorm:"column:position;default:''" json:"position"`       // 部门ID
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`     // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`     // 更新时间
}

func (Organize) TableName() string {
	return "organize"
}
