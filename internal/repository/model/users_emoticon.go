package model

import "time"

type UsersEmoticon struct {
	Id          int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 表情包收藏ID
	UserId      int       `gorm:"column:user_id;NOT NULL" json:"user_id"`           // 用户ID
	EmoticonIds string    `gorm:"column:emoticon_ids;NOT NULL" json:"emoticon_ids"` // 表情包ID
	CreatedAt   time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`     // 创建时间
}
