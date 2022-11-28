package model

import "time"

type ArticleTag struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 文章分类ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 用户ID
	TagName   string    `gorm:"column:tag_name;NOT NULL" json:"tag_name"`         // 标签名
	Sort      int       `gorm:"column:sort;default:0;NOT NULL" json:"sort"`       // 排序
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`     // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`     // 更新时间
}

func (ArticleTag) TableName() string {
	return "article_tag"
}

type TagItem struct {
	Id      int    `json:"id"`       // 文章分类ID
	TagName string `json:"tag_name"` // 标签名
	Count   int    `json:"count"`    // 排序
}
