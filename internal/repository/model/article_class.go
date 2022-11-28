package model

import "time"

type ArticleClass struct {
	Id        int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`         // 文章分类ID
	UserId    int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`       // 用户ID
	ClassName string    `gorm:"column:class_name;NOT NULL" json:"class_name"`           // 分类名
	Sort      int       `gorm:"column:sort;default:0;NOT NULL" json:"sort"`             // 排序
	IsDefault int       `gorm:"column:is_default;default:0;NOT NULL" json:"is_default"` // 默认分类[0:否;1:是；]
	CreatedAt time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`           // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`           // 更新时间
}

func (ArticleClass) TableName() string {
	return "article_class"
}

type ArticleClassItem struct {
	Id        int    `json:"id"`         // 文章分类ID
	ClassName string `json:"class_name"` // 分类名
	IsDefault int    `json:"is_default"` // 默认分类1:是 0:不是
	Count     int    `json:"count"`      // 分类名
}
