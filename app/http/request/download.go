package request

type DownloadChatFileRequest struct {
	CrId int `form:"cr_id" json:"cr_id" binding:"required,min=1"`
}
