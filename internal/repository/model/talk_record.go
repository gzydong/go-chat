package model

import "time"

type TalkRecord struct {
	MsgId      string    `gorm:"column:msg_id;" json:"msg_id"`           // 消息唯一ID
	Sequence   int64     `gorm:"column:sequence;" json:"sequence"`       // 消息时序ID
	TalkType   int       `gorm:"column:talk_type;" json:"talk_type"`     // 对话类型[1:私信;2:群聊;]
	MsgType    int       `gorm:"column:msg_type;" json:"msg_type"`       // 消息类型
	UserId     int       `gorm:"column:user_id;" json:"user_id"`         // 发送者ID[0:系统用户;]
	ReceiverId int       `gorm:"column:receiver_id;" json:"receiver_id"` // 接收者ID（用户ID 或 群ID）
	IsRevoke   int       `gorm:"column:is_revoked;" json:"is_revoked"`   // 是否撤回消息[0:否;1:是;]
	QuoteId    string    `gorm:"column:quote_id;" json:"quote_id"`       // 引用消息ID
	Extra      string    `gorm:"column:extra;" json:"extra"`             // 扩展信信息
	CreatedAt  time.Time `gorm:"column:created_at;" json:"created_at"`   // 创建时间
}
