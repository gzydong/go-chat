package request

type UploadFileStreamRequest struct {
	Stream string `form:"stream"`
}

type UploadInitiateMultipartRequest struct {
	FileName string `form:"file_name" binding:"required"`
	FileSize int64  `form:"file_size" binding:"required"`
}
