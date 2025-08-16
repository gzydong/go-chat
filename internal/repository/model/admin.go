package model

import (
	"time"

	"go-chat/internal/pkg/encrypt"
)

const (
	AdminStatusNormal   = 1
	AdminStatusDisabled = 2
)

type Admin struct {
	Id          int       `gorm:"column:id" db:"id" json:"id" form:"id"`                                             // 用户ID
	Email       string    `gorm:"column:email" db:"email" json:"email" form:"email"`                                 // 邮箱
	Username    string    `gorm:"column:username" db:"username" json:"username" form:"username"`                     // 用户昵称
	Password    string    `gorm:"column:password" db:"password" json:"password" form:"password"`                     // 用户密码
	Avatar      string    `gorm:"column:avatar" db:"avatar" json:"avatar" form:"avatar"`                             // 用户头像
	Gender      int8      `gorm:"column:gender" db:"gender" json:"gender" form:"gender"`                             // 用户性别[1:男;2:女;3:未知]
	Mobile      string    `gorm:"column:mobile" db:"mobile" json:"mobile" form:"mobile"`                             // 手机号
	Motto       string    `gorm:"column:motto" db:"motto" json:"motto" form:"motto"`                                 // 座右铭
	Address     string    `gorm:"column:address" db:"address" json:"address" form:"address"`                         // 联系地址
	Status      int       `gorm:"column:status" db:"status" json:"status" form:"status"`                             // 状态 1正常 2停用
	LastLoginAt time.Time `gorm:"column:last_login_at" db:"last_login_at" json:"last_login_at" form:"last_login_at"` // 注册时间
	CreatedAt   time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"`             // 注册时间
	UpdatedAt   time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" form:"updated_at"`             // 更新时间
}

func (Admin) TableName() string {
	return "admin"
}

func (a Admin) IsDisabled() bool {
	return a.Status == AdminStatusDisabled
}

func (a Admin) VerifyPassword(password string) bool {
	return encrypt.VerifyPassword(a.Password, password)
}
