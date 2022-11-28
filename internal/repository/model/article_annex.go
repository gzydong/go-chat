package model

import (
	"database/sql"
	"time"
)

type ArticleAnnex struct {
	Id           int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`         // 文件ID
	UserId       int          `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`       // 上传文件的用户ID
	ArticleId    int          `gorm:"column:article_id;default:1;NOT NULL" json:"article_id"` // 笔记ID
	Drive        int          `gorm:"column:drive;default:1;NOT NULL" json:"drive"`           // 文件驱动[1:local;2:cos;]
	Suffix       string       `gorm:"column:suffix;NOT NULL" json:"suffix"`                   // 文件后缀名
	Size         int          `gorm:"column:size;default:0;NOT NULL" json:"size"`             // 文件大小
	Path         string       `gorm:"column:path;NOT NULL" json:"path"`                       // 文件地址（相对地址）
	OriginalName string       `gorm:"column:original_name;NOT NULL" json:"original_name"`     // 原文件名
	Status       int          `gorm:"column:status;default:1;NOT NULL" json:"status"`         // 附件状态[1:正常;2:已删除;]
	CreatedAt    time.Time    `gorm:"column:created_at;NOT NULL" json:"created_at"`           // 创建时间
	UpdatedAt    time.Time    `gorm:"column:updated_at;NOT NULL" json:"updated_at"`           // 更新时间
	DeletedAt    sql.NullTime `gorm:"column:deleted_at" json:"deleted_at"`                    // 删除时间
}

func (ArticleAnnex) TableName() string {
	return "article_annex"
}

type RecoverAnnexItem struct {
	Id           int       `json:"id"`            // 文件ID
	ArticleId    int       `json:"article_id"`    // 笔记ID
	Title        string    `json:"title"`         // 原文件名
	OriginalName string    `json:"original_name"` // 原文件名
	DeletedAt    time.Time `json:"deleted_at"`    // 附件删除时间
}
