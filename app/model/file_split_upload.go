package model

import "time"

type FileSplitUpload struct {
	Id           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`  // 临时文件ID
	FileType     int       `gorm:"column:file_type;default:1" json:"file_type"`     // 1:合并文件  2:拆分文件
	UserId       int       `gorm:"column:user_id;default:0" json:"user_id"`         // 上传的用户ID
	HashName     string    `gorm:"column:hash_name" json:"hash_name"`               // 临时文件hash名
	OriginalName string    `gorm:"column:original_name" json:"original_name"`       // 原文件名
	SplitIndex   int       `gorm:"column:split_index;default:0" json:"split_index"` // 当前索引块
	SplitNum     int       `gorm:"column:split_num;default:0" json:"split_num"`     // 总上传索引块
	SaveDir      string    `gorm:"column:save_dir" json:"save_dir"`                 // 文件的临时保存路径
	FileExt      string    `gorm:"column:file_ext" json:"file_ext"`                 // 文件后缀名
	FileSize     int64     `gorm:"column:file_size" json:"file_size"`               // 临时文件大小
	IsDelete     int       `gorm:"column:is_delete;default:0" json:"is_delete"`     // 文件是否已被删除(1:是 0:否)
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`             // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`             // 更新时间
}
