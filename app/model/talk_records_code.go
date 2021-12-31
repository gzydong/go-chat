package model

import "time"

type TalkRecordsCode struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`       // 代码块ID
	RecordId  int       `gorm:"column:record_id;default:0;NOT NULL" json:"record_id"` // 消息记录ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`     // 用户ID
	Lang      string    `gorm:"column:lang;NOT NULL" json:"lang"`                     // 语言类型(如：php,java,python)
	Code      string    `gorm:"column:code;NOT NULL" json:"code"`                     // 代码内容
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`         // 创建时间
}
