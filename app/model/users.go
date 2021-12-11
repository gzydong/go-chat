package model

import "time"

type Users struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 用户ID
	Mobile    string    `gorm:"column:mobile" json:"mobile"`                    // 手机号
	Nickname  string    `gorm:"column:nickname" json:"nickname"`                // 用户昵称
	Avatar    string    `gorm:"column:avatar" json:"avatar"`                    // 用户头像地址
	Gender    int       `gorm:"column:gender;default:0" json:"gender"`          // 用户性别  0:未知  1:男   2:女
	Password  string    `gorm:"column:password;NOT NULL" json:"password"`       // 用户密码
	Motto     string    `gorm:"column:motto" json:"motto"`                      // 用户座右铭
	Email     string    `gorm:"column:email" json:"email"`                      // 用户邮箱
	IsRobot   int       `gorm:"column:is_robot;default:0" json:"is_robot"`      // 是否机器人[0:否;1:是;]
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`            // 注册时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`            // 更新时间
}

func (m *Users) TableName() string {
	return "users"
}
