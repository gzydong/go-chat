package model

import "time"

type TalkRecordsFile struct {
	Id           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 文件ID
	RecordId     int       `gorm:"column:record_id;default:0" json:"record_id"`      // 消息记录ID
	UserId       int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 上传文件的用户ID
	Source       int       `gorm:"column:source;default:1" json:"-"`                 // 文件来源[1:用户上传;2:表情包;]
	Type         int       `gorm:"column:type;default:1" json:"type"`                // 文件类型[1:图片;2:音频文件;3:视频文件;4:其它文件;]
	Drive        int       `gorm:"column:drive;default:1" json:"-"`                  // 文件保存方式[1:local;2:cos;]
	OriginalName string    `gorm:"column:original_name" json:"original_name"`        // 原文件名
	Suffix       string    `gorm:"column:suffix" json:"suffix"`                      // 文件后缀名
	Size         int       `gorm:"column:size;default:0" json:"size"`                // 文件大小
	Path         string    `gorm:"column:path" json:"path"`                          // 文件地址
	Url          string    `gorm:"column:url;NOT NULL" json:"url"`                   // 网络连接
	IsDelete     int       `gorm:"column:is_delete;default:0" json:"-"`              // 文件是否已删除 （0:否 1:已删除）
	CreatedAt    time.Time `gorm:"column:created_at" json:"-"`                       // 创建时间
}
