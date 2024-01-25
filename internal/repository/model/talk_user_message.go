package model

import "time"

type TalkUserMessage struct {
	Id        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 聊天记录ID
	MsgId     string    `gorm:"column:msg_id;" json:"msg_id"`                   // 消息ID
	OrgMsgId  string    `gorm:"column:org_msg_id;" json:"org_msg_id"`           // 原消息ID
	Sequence  int64     `gorm:"column:sequence;" json:"sequence"`               // 消息时序ID（消息排序）
	MsgType   int       `gorm:"column:msg_type;" json:"msg_type"`               // 消息类型
	UserId    int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	ToFromId  int       `gorm:"column:to_from_id;" json:"to_from_id"`           // 接受者ID
	FromId    int       `gorm:"column:from_id;" json:"from_id"`                 // 消息发送者ID
	IsRevoked int       `gorm:"column:is_revoked;" json:"is_revoked"`           // 是否撤回[1:否;2:是;]
	IsDeleted int       `gorm:"column:is_deleted;" json:"is_deleted"`           // 是否删除[1:否;2:是;]
	Extra     string    `gorm:"column:extra;" json:"extra"`                     // 消息扩展字段
	Quote     string    `gorm:"column:quote;" json:"quote"`                     // 引用消息ID
	SendTime  time.Time `gorm:"column:send_time;" json:"send_time"`             // 发送时间
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 更新时间
}

func (TalkUserMessage) TableName() string {
	return "talk_user_message"
}
