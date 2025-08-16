package model

import (
	"time"
)

type SysMenu struct {
	Id        int32     `gorm:"column:id" db:"id" json:"id" form:"id"`                                 // 用户ID
	ParentId  int32     `gorm:"column:parent_id" db:"parent_id" json:"parent_id" form:"parent_id"`     // 用户昵称
	Name      string    `gorm:"column:name" db:"name" json:"name" form:"name"`                         // 用户密码
	MenuType  int32     `gorm:"column:menu_type" db:"menu_type" json:"menu_type" form:"menu_type"`     // 用户头像
	Icon      string    `gorm:"column:icon" db:"icon" json:"icon" form:"icon"`                         // 手机号
	Path      string    `gorm:"column:path" db:"path" json:"path" form:"path"`                         // 邮箱
	Sort      int32     `gorm:"column:sort" db:"sort" json:"sort" form:"sort"`                         // 座右铭
	Hidden    string    `gorm:"column:hidden" db:"hidden" json:"hidden" form:"hidden"`                 // 状态 1正常 2停用
	UseLayout string    `gorm:"column:use_layout" db:"use_layout" json:"use_layout" form:"use_layout"` // 状态 Y:使用 2:未使用
	AuthCode  string    `gorm:"column:auth_code" db:"auth_code" json:"auth_code" form:"auth_code"`     // 状态 Y:使用 2:未使用
	Status    int32     `gorm:"column:status" db:"status" json:"status" form:"status"`                 // 状态 1正常 2停用
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"` // 注册时间
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" form:"updated_at"` // 更新时间
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
