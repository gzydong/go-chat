package model

type TalkRecordsFile struct {
	ID           int    `json:"id"`
	RecordId     string `json:"record_id"`
	UserId       string `json:"user_id"`
	FileSource   string `json:"file_source"`
	FileType     string `json:"file_type"`
	SaveType     string `json:"save_type"`
	OriginalName string `json:"original_name"`
	FileSuffix   string `json:"file_suffix"`
	FileSize     string `json:"file_size"`
	SaveDir      string `json:"save_dir"`
	IsDelete     string `json:"is_delete"`
	CreatedAt    string `json:"created_at"`
}
