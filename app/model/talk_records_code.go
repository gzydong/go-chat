package model

import "time"

type TalkRecordsCode struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"-"` // 代码块ID
	RecordId  int       `gorm:"column:record_id;default:0" json:"-"`           // 消息记录ID
	UserId    int       `gorm:"column:user_id;default:0" json:"-"`             // 用户ID
	CodeLang  string    `gorm:"column:code_lang" json:"code_lang"`             // 代码片段类型(如：php,java,python)
	Code      string    `gorm:"column:code" json:"code"`                       // 代码片段内容
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`           // 创建时间
}
