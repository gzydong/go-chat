package model

import "time"

type TalkRecordsRead struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`           // 自增ID
	MsgId      string    `gorm:"column:msg_id;NOT NULL" json:"msg_id"`                     // 消息ID
	UserId     int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`         // 用户ID
	ReceiverId int       `gorm:"column:receiver_id;default:0;NOT NULL" json:"receiver_id"` // 接受者ID
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`             // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`             // 更新时间
}

func (t TalkRecordsRead) TableName() string {
	return "talk_records_read"
}
