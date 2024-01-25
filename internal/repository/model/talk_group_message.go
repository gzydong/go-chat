package model

import "time"

type TalkGroupMessage struct {
	Id        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 聊天记录ID
	MsgId     string    `gorm:"column:msg_id;" json:"msg_id"`                   // 消息ID
	Sequence  int64     `gorm:"column:sequence;" json:"sequence"`               // 消息时序ID（消息排序）
	MsgType   int       `gorm:"column:msg_type;" json:"msg_type"`               // 消息类型
	GroupId   int       `gorm:"column:group_id;" json:"group_id"`               // 群组ID
	FromId    int       `gorm:"column:from_id;" json:"from_id"`                 // 消息发送者ID
	IsRevoked int       `gorm:"column:is_revoked;" json:"is_revoked"`           // 是否撤回[1:否;2:是;]
	Extra     string    `gorm:"column:extra;" json:"extra"`                     // 消息扩展字段
	Quote     string    `gorm:"column:quote;" json:"quote"`                     // 引用消息
	SendTime  time.Time `gorm:"column:send_time;" json:"send_time"`             // 发送时间
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (TalkGroupMessage) TableName() string {
	return "talk_group_message"
}
