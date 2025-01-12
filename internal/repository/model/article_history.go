package model

import (
	"time"
)

type ArticleHistory struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 文章ID
	UserId    int       `gorm:"column:user_id;" json:"user_id"`                 // 用户ID
	ArticleId int       `gorm:"column:article_id;" json:"article_id"`           // 笔记ID
	Content   string    `gorm:"column:content;" json:"content"`                 // Markdown 内容
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at"`           // 创建时间
}

func (ArticleHistory) TableName() string {
	return "article_history"
}
