package model

import (
	"database/sql"
	"time"
)

type ArticleAnnex struct {
	Id           int          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 文件ID
	UserId       int          `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 上传文件的用户ID
	ArticleId    int          `gorm:"column:article_id;default:1" json:"article_id"`    // 笔记ID
	FileSuffix   string       `gorm:"column:file_suffix" json:"file_suffix"`            // 文件后缀名
	FileSize     int          `gorm:"column:file_size;default:0" json:"file_size"`      // 文件大小（单位字节）
	SaveDir      string       `gorm:"column:save_dir" json:"save_dir"`                  // 文件保存地址（相对地址）
	OriginalName string       `gorm:"column:original_name" json:"original_name"`        // 原文件名
	Status       int          `gorm:"column:status;default:1" json:"status"`            // 附件状态 1:正常 2:已删除
	CreatedAt    time.Time    `gorm:"column:created_at" json:"created_at"`              // 附件上传时间
	DeletedAt    sql.NullTime `gorm:"column:deleted_at" json:"deleted_at"`              // 附件删除时间
}
