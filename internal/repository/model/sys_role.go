package model

import (
	"time"
)

type SysRole struct {
	Id        int       `gorm:"column:id" db:"id" json:"id"`
	RoleName  string    `gorm:"column:role_name" db:"role_name" json:"role_name"` // 用户昵称
	Status    int       `gorm:"column:status" db:"status" json:"status"`
	Explain   string    `gorm:"column:explain" db:"explain" json:"explain"`
	CreatedAt time.Time `gorm:"column:created_at;" db:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" db:"updated_at" json:"updated_at"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
