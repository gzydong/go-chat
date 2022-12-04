package model

import "time"

type Users struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`     // 用户ID
	Mobile    string    `gorm:"column:mobile;NOT NULL" json:"mobile"`               // 手机号
	Nickname  string    `gorm:"column:nickname;NOT NULL" json:"nickname"`           // 用户昵称
	Avatar    string    `gorm:"column:avatar;NOT NULL" json:"avatar"`               // 用户头像地址
	Gender    int       `gorm:"column:gender;default:0;NOT NULL" json:"gender"`     // 用户性别  0:未知  1:男   2:女
	Password  string    `gorm:"column:password;NOT NULL" json:"-"`                  // 用户密码
	Motto     string    `gorm:"column:motto;NOT NULL" json:"motto"`                 // 用户座右铭
	Email     string    `gorm:"column:email;NOT NULL" json:"email"`                 // 用户邮箱
	Birthday  string    `gorm:"column:birthday;NOT NULL" json:"birthday"`           // 生日
	IsRobot   int       `gorm:"column:is_robot;default:0;NOT NULL" json:"is_robot"` // 是否机器人[0:否;1:是;]
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`       // 注册时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`       // 更新时间
}

func (Users) TableName() string {
	return "users"
}
