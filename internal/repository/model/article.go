package model

import (
	"database/sql"
	"time"
)

type Article struct {
	Id         int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`           // 文章ID
	UserId     int          `gorm:"column:user_id;NOT NULL" json:"user_id"`                   // 用户ID
	ClassId    int          `gorm:"column:class_id;default:0;NOT NULL" json:"class_id"`       // 分类ID
	TagsId     string       `gorm:"column:tags_id;NOT NULL" json:"tags_id"`                   // 笔记关联标签
	Title      string       `gorm:"column:title;NOT NULL" json:"title"`                       // 文章标题
	Abstract   string       `gorm:"column:abstract;NOT NULL" json:"abstract"`                 // 文章摘要
	Image      string       `gorm:"column:image;NOT NULL" json:"image"`                       // 文章首图
	IsAsterisk int          `gorm:"column:is_asterisk;default:0;NOT NULL" json:"is_asterisk"` // 是否星标文章[0:否;1:是;]
	Status     int          `gorm:"column:status;default:1;NOT NULL" json:"status"`           // 笔记状态[1:正常;2:已删除;]
	MdContent  string       `gorm:"column:md_content;NOT NULL" json:"md_content"`             // Markdown 内容
	CreatedAt  time.Time    `gorm:"column:created_at;NOT NULL" json:"created_at"`             // 创建时间
	UpdatedAt  time.Time    `gorm:"column:updated_at;NOT NULL" json:"updated_at"`             // 更新时间
	DeletedAt  sql.NullTime `gorm:"column:deleted_at" json:"deleted_at"`                      // 删除时间
}

func (Article) TableName() string {
	return "article"
}

type ArticleListItem struct {
	Id         int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`           // 文章ID
	UserId     int       `gorm:"column:user_id;NOT NULL" json:"user_id"`                   // 用户ID
	ClassId    int       `gorm:"column:class_id;default:0;NOT NULL" json:"class_id"`       // 分类ID
	TagsId     string    `gorm:"column:tags_id;NOT NULL" json:"tags_id"`                   // 笔记关联标签
	Title      string    `gorm:"column:title;NOT NULL" json:"title"`                       // 文章标题
	Abstract   string    `gorm:"column:abstract;NOT NULL" json:"abstract"`                 // 文章摘要
	Image      string    `gorm:"column:image;NOT NULL" json:"image"`                       // 文章首图
	IsAsterisk int       `gorm:"column:is_asterisk;default:0;NOT NULL" json:"is_asterisk"` // 是否星标文章[0:否;1:是;]
	Status     int       `gorm:"column:status;default:1;NOT NULL" json:"status"`           // 笔记状态[1:正常;2:已删除;]
	CreatedAt  time.Time `gorm:"column:created_at;NOT NULL" json:"created_at"`             // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;NOT NULL" json:"updated_at"`             // 更新时间
	ClassName  string    `gorm:"column:class_name;NOT NULL" json:"class_name"`             // 分类名
}

type ArticleDetailInfo struct {
	Id         int       `json:"id"`          // 文章ID
	UserId     int       `json:"user_id"`     // 用户ID
	ClassId    int       `json:"class_id"`    // 分类ID
	TagsId     string    `json:"tags_id"`     // 笔记关联标签
	Title      string    `json:"title"`       // 文章标题
	Abstract   string    `json:"abstract"`    // 文章摘要
	Image      string    `json:"image"`       // 文章首图
	IsAsterisk int       `json:"is_asterisk"` // 是否星标文章(0:否  1:是)
	Status     int       `json:"status"`      // 笔记状态 1:正常 2:已删除
	CreatedAt  time.Time `json:"created_at"`  // 添加时间
	UpdatedAt  time.Time `json:"updated_at"`  // 最后一次更新时间
	MdContent  string    `json:"md_content"`  // Markdown 内容
}
