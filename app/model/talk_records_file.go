package model

import "time"

type TalkRecordsFile struct {
	ID           int       `json:"id" grom:"comment:文件消息ID"`
	RecordId     int       `json:"record_id" grom:"comment:聊天记录ID"`
	UserId       int       `json:"user_id" grom:"comment:用户ID"`
	FileSource   int       `json:"file_source" grom:"comment:文件上传来源"`
	FileType     int       `json:"file_type" grom:"comment:文件类型"`
	SaveType     int       `json:"save_type" grom:"comment:文件保存类型"`
	OriginalName string    `json:"original_name" grom:"comment:文件原始名称"`
	FileSuffix   string    `json:"file_suffix" grom:"comment:文件后缀名"`
	FileSize     int       `json:"file_size" grom:"comment:文件大小"`
	SaveDir      string    `json:"save_dir" grom:"comment:文件保存路径"`
	IsDelete     int       `json:"-" grom:"comment:是否已删除"`
	CreatedAt    time.Time `json:"-" grom:"comment:上传时间"`
}
