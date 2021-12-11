package model

import "time"

type TalkRecordsFile struct {
	Id           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`   // 文件ID
	RecordId     int       `gorm:"column:record_id;default:0" json:"record_id"`      // 消息记录ID
	UserId       int       `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"` // 上传文件的用户ID
	FileSource   int       `gorm:"column:file_source;default:1" json:"file_source"`  // 文件来源（1:用户上传 2:表情包）
	FileType     int       `gorm:"column:file_type;default:1" json:"file_type"`      // 文件类型（1:图片 2:音频文件 3:视频文件 4:其它文件 ）
	SaveType     int       `gorm:"column:save_type;default:0" json:"save_type"`      // 文件保存方式（0:本地 1:第三方[阿里OOS、七牛云] ）
	OriginalName string    `gorm:"column:original_name" json:"original_name"`        // 原文件名
	FileSuffix   string    `gorm:"column:file_suffix" json:"file_suffix"`            // 文件后缀名
	FileSize     int       `gorm:"column:file_size;default:0" json:"file_size"`      // 文件大小（单位字节）
	SaveDir      string    `gorm:"column:save_dir" json:"save_dir"`                  // 文件保存地址（相对地址/第三方网络地址）
	IsDelete     int       `gorm:"column:is_delete;default:0" json:"-"`              // 文件是否已删除 （0:否 1:已删除）
	CreatedAt    time.Time `gorm:"column:created_at" json:"-"`                       // 创建时间
}
