package model

import "time"

type RobotInstallUser struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`     // 自增ID
	RobotId   int       `gorm:"column:robot_id;default:0;NOT NULL" json:"robot_id"` // 机器人ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`   // 关联用户ID
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`       // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`       // 更新时间
}

func (RobotInstallUser) TableName() string {
	return "robot_install_user"
}
