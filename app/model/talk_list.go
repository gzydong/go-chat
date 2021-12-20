package model

import "time"

type TalkList struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`  // 聊天列表ID
	TalkType   int       `gorm:"column:talk_type;default:1" json:"talk_type"`     // 聊天类型[1:私信;2:群聊;]
	UserId     int       `gorm:"column:user_id;default:0" json:"user_id"`         // 用户ID
	ReceiverId int       `gorm:"column:receiver_id;default:0" json:"receiver_id"` // 接收者ID（用户ID 或 群ID）
	IsTop      int       `gorm:"column:is_top;default:0" json:"is_top"`           // 是否置顶[0:否;1:是;]
	IsDisturb  int       `gorm:"column:is_disturb;default:0" json:"is_disturb"`   // 消息免打扰[0:否;1:是;]
	IsDelete   int       `gorm:"column:is_delete;default:0" json:"is_delete"`     // 是否删除[0:否;1:是;]
	IsRobot    int       `gorm:"column:is_robot;default:0" json:"is_robot"`       // 是否机器人[0:否;1:是;]
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`             // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`             // 更新时间
}

type SearchTalkList struct {
	Id          int       `json:"id" `
	TalkType    int       `json:"talk_type" `
	ReceiverId  int       `json:"receiver_id" `
	IsDelete    int       `json:"is_delete"`
	IsTop       int       `json:"is_top"`
	IsRobot     int       `json:"is_robot"`
	IsDisturb   int       `json:"is_disturb"`
	UserAvatar  string    `json:"user_avatar"`
	Nickname    string    `json:"nickname"`
	GroupName   string    `json:"group_name"`
	GroupAvatar string    `json:"group_avatar"`
	UpdatedAt   time.Time `json:"updated_at"`
}
