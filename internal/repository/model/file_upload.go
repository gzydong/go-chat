package model

import (
	"time"
)

type FileUpload struct {
	Id           int       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"` // 临时文件ID
	Type         int       `gorm:"column:type;" json:"type"`                       // 文件属性[1:合并文件;2:拆分文件]
	Drive        int       `gorm:"column:drive;" json:"drive"`                     // 驱动类型[1:local;2:minio;]
	UploadId     string    `gorm:"column:upload_id;" json:"upload_id"`             // 临时文件hash名
	UserId       int       `gorm:"column:user_id;" json:"user_id"`                 // 上传的用户ID
	OriginalName string    `gorm:"column:original_name;" json:"original_name"`     // 原文件名
	SplitIndex   int       `gorm:"column:split_index;" json:"split_index"`         // 当前索引块
	SplitNum     int       `gorm:"column:split_num;" json:"split_num"`             // 总上传索引块
	Path         string    `gorm:"column:path;" json:"path"`                       // 临时保存路径
	FileExt      string    `gorm:"column:file_ext;" json:"file_ext"`               // 文件后缀名
	FileSize     int64     `gorm:"column:file_size;" json:"file_size"`             // 文件大小
	IsDelete     int       `gorm:"column:is_delete;" json:"is_delete"`             // 文件是否删除[0:否;1:是;]
	Attr         string    `gorm:"column:attr;" json:"attr"`                       // 额外参数json
	CreatedAt    time.Time `gorm:"column:created_at;" json:"created_at"`           // 更新时间
	UpdatedAt    time.Time `gorm:"column:updated_at;" json:"updated_at"`           // 创建时间
}

func (FileUpload) TableName() string {
	return "file_upload"
}
