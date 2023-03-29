package model

import "time"

const (
	AdminStatusNormal     = 1
	AdminStatusDeactivate = 2
)

type Admin struct {
	Id        int       `gorm:"column:id" db:"id" json:"id" form:"id"`                                 // 用户ID
	Username  string    `gorm:"column:username" db:"username" json:"username" form:"username"`         // 用户昵称
	Password  string    `gorm:"column:password" db:"password" json:"password" form:"password"`         // 用户密码
	Avatar    string    `gorm:"column:avatar" db:"avatar" json:"avatar" form:"avatar"`                 // 用户头像
	Gender    int8      `gorm:"column:gender" db:"gender" json:"gender" form:"gender"`                 // 用户性别[0:未知;1:男 ;2:女;]
	Mobile    string    `gorm:"column:mobile" db:"mobile" json:"mobile" form:"mobile"`                 // 手机号
	Email     string    `gorm:"column:email" db:"email" json:"email" form:"email"`                     // 用户邮箱
	Motto     string    `gorm:"column:motto" db:"motto" json:"motto" form:"motto"`                     // 用户座右铭
	Birthday  string    `gorm:"column:birthday" db:"birthday" json:"birthday" form:"birthday"`         // 生日
	Status    int8      `gorm:"column:status" db:"status" json:"status" form:"status"`                 // 状态 1正常 2停用
	CreatedAt time.Time `gorm:"column:created_at" db:"created_at" json:"created_at" form:"created_at"` // 注册时间
	UpdatedAt time.Time `gorm:"column:updated_at" db:"updated_at" json:"updated_at" form:"updated_at"` // 更新时间
}

func (Admin) TableName() string {
	return "admin"
}
