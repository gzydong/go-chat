package web

type UploadInitiateMultipartRequest struct {
	FileName string `form:"file_name" json:"file_name" binding:"required"`
	FileSize int64  `form:"file_size" json:"file_size" binding:"required"`
}

type UploadMultipartRequest struct {
	UploadId   string `form:"upload_id" json:"upload_id" binding:"required"`
	SplitIndex int    `form:"split_index" json:"split_index" binding:"min=0"`
	SplitNum   int    `form:"split_num" json:"split_num" binding:"required,min=1"`
}
