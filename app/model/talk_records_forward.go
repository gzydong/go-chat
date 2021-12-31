package model

import "time"

type TalkRecordsForward struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`       // 合并转发ID
	RecordId  int       `gorm:"column:record_id;default:0;NOT NULL" json:"record_id"` // 消息记录ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`     // 转发用户ID
	RecordsId string    `gorm:"column:records_id;NOT NULL" json:"records_id"`         // 转发的聊天记录ID （多个用 , 拼接），最多只能30条记录信息
	Text      string    `gorm:"column:text;NOT NULL" json:"text"`                     // 记录快照（避免后端再次查询数据）
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`         // 转发时间
}
