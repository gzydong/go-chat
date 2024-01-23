package model

import "time"

const (
	TalkRecordTalkTypePrivate = 1
	TalkRecordTalkTypeGroup   = 2
)

type TalkRecords struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`           // 自增ID
	MsgId      string    `gorm:"column:msg_id;NOT NULL" json:"msg_id"`                     // 消息唯一ID
	Sequence   int64     `gorm:"column:sequence;default:0;NOT NULL" json:"sequence"`       // 消息时序ID
	TalkType   int       `gorm:"column:talk_type;default:1;NOT NULL" json:"talk_type"`     // 对话类型[1:私信;2:群聊;]
	MsgType    int       `gorm:"column:msg_type;default:0;NOT NULL" json:"msg_type"`       // 消息类型
	UserId     int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`         // 发送者ID[0:系统用户;]
	ReceiverId int       `gorm:"column:receiver_id;default:0;NOT NULL" json:"receiver_id"` // 接收者ID（用户ID 或 群ID）
	IsRevoke   int       `gorm:"column:is_revoke;default:0;NOT NULL" json:"is_revoke"`     // 是否撤回消息[0:否;1:是;]
	QuoteId    string    `gorm:"column:quote_id;NOT NULL" json:"quote_id"`                 // 引用消息ID
	Extra      string    `gorm:"column:extra;default:{}" json:"extra"`                     // 扩展信信息
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`             // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`             // 更新时间
}

func (t TalkRecords) TableName() string {
	return "talk_records"
}
