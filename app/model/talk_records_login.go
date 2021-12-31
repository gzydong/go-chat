package model

import "time"

type TalkRecordsLogin struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`       // 登录ID
	RecordId  int       `gorm:"column:record_id;default:0;NOT NULL" json:"record_id"` // 消息记录ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`     // 用户ID
	Ip        string    `gorm:"column:ip;NOT NULL" json:"ip"`                         // IP地址
	Platform  string    `gorm:"column:platform;NOT NULL" json:"platform"`             // 登录平台[h5,ios,windows,mac,web]
	Agent     string    `gorm:"column:agent;NOT NULL" json:"agent"`                   // 设备信息
	Address   string    `gorm:"column:address;NOT NULL" json:"address"`               // IP所在地
	Reason    string    `gorm:"column:reason;NOT NULL" json:"reason"`                 // 异常提示
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`         // 创建时间
}
