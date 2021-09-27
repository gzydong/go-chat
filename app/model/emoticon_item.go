package model

type EmoticonItem struct {
	ID         int    `json:"id" grom:"comment:表情包详情ID"`
	EmoticonId string `json:"emoticon_id" grom:"comment:表情分组ID"`
	UserId     string `json:"user_id" grom:"comment:用户ID"`
	Describe   string `json:"describe" grom:"comment:表情描述"`
	Url        string `json:"url" grom:"comment:表情链接"`
	FileSuffix string `json:"file_suffix" grom:"comment:文件前缀"`
	FileSize   string `json:"file_size" grom:"comment:表情包文件大小"`
	CreatedAt  string `json:"created_at" grom:"comment:创建时间"`
	UpdatedAt  string `json:"updated_at" grom:"comment:更新时间"`
}
