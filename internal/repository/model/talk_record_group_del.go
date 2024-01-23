package model

import "time"

type TalkRecordGroupDel struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 用户ID
	MsgId     string    `gorm:"column:msg_id;default:'';NOT NULL" json:"msg_id"`  // 聊天记录ID
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`     // 创建时间
}

func (t TalkRecordGroupDel) TableName() string {
	return "talk_record_group_del"
}
