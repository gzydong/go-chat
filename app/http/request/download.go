package request

type DownloadChatFileRequest struct {
	RecordId int `form:"cr_id" json:"cr_id" binding:"required,min=1"`
}
