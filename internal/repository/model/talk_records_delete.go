package model

import "time"

type TalkRecordsDelete struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	MsgId     string    `gorm:"column:msg_id;default:'';NOT NULL" json:"msg_id"`  // 聊天记录ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 用户ID
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`     // 创建时间
}

func (t TalkRecordsDelete) TableName() string {
	return "talk_records_delete"
}
