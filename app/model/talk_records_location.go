package model

import "time"

type TalkRecordsLocation struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`     // 自增ID
	RecordId  int       `gorm:"column:record_id;default:0" json:"record_id"`        // 消息记录ID
	UserId    int       `gorm:"column:user_id;default:0" json:"user_id"`            // 用户ID
	Longitude string    `gorm:"column:longitude;default:0.000000" json:"longitude"` // 经度
	Latitude  string    `gorm:"column:latitude;default:0.000000" json:"latitude"`   // 纬度
	CreatedAt time.Time `gorm:"column:created_at" json:"-"`                         // 登录时间
}
