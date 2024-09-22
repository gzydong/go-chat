package model

import "time"

type TalkSession struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 聊天列表ID
	TalkMode  int       `gorm:"column:talk_mode;" json:"talk_mode"`             // 聊天类型[1:私信;2:群聊;]
	UserId    int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	ToFromId  int       `gorm:"column:to_from_id;" json:"to_from_id"`           // 接收者ID（用户ID 或 群ID）
	IsTop     int       `gorm:"column:is_top;" json:"is_top"`                   // 是否置顶[1:否;2:是;]
	IsDisturb int       `gorm:"column:is_disturb;" json:"is_disturb"`           // 消息免打扰[1:否;2:是;]
	IsDelete  int       `gorm:"column:is_delete;" json:"is_delete"`             // 是否删除[1:否;2:是;]
	IsRobot   int       `gorm:"column:is_robot;" json:"is_robot"`               // 是否机器人[1:否;2:是;]
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (TalkSession) TableName() string {
	return "talk_session"
}

type SearchTalkSession struct {
	Id          int       `json:"id"`
	TalkMode    int       `json:"talk_mode"`
	ToFromId    int       `json:"to_from_id" `
	IsDelete    int       `json:"is_delete"`
	IsTop       int       `json:"is_top"`
	IsRobot     int       `json:"is_robot"`
	IsDisturb   int       `json:"is_disturb"`
	Avatar      string    `json:"avatar"`
	Nickname    string    `json:"nickname"`
	GroupName   string    `json:"group_name"`
	GroupAvatar string    `json:"group_avatar"`
	UpdatedAt   time.Time `json:"updated_at"`
}
