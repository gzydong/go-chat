package model

import (
	"time"
)

const (
	UsersGenderDefault = 3
)

type Users struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 用户ID
	Mobile    string    `gorm:"column:mobile;" json:"mobile"`                   // 手机号
	Nickname  string    `gorm:"column:nickname;" json:"nickname"`               // 用户昵称
	Avatar    string    `gorm:"column:avatar;" json:"avatar"`                   // 用户头像地址
	Gender    int       `gorm:"column:gender;" json:"gender"`                   // 用户性别 1:男 2:女 3:未知
	Password  string    `gorm:"column:password;" json:"-"`                      // 用户密码
	Motto     string    `gorm:"column:motto;" json:"motto"`                     // 用户座右铭
	Email     string    `gorm:"column:email;" json:"email"`                     // 用户邮箱
	Birthday  string    `gorm:"column:birthday;" json:"birthday"`               // 生日
	IsRobot   int       `gorm:"column:is_robot;" json:"is_robot"`               // 是否机器人[1:否;2:是;]
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 注册时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (u Users) TableName() string {
	return "users"
}

func (u Users) TablePrimaryId() string {
	return "id"
}

func (u Users) TablePrimaryIdValue() int {
	return u.Id
}
