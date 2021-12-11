package model

import "time"

type TalkRecordsDelete struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	RecordId  int       `gorm:"column:record_id;default:0" json:"record_id"` // 聊天记录ID
	UserId    int       `gorm:"column:user_id;default:0" json:"user_id"`     // 用户ID
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`         // 删除时间
}
