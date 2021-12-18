package request

type UploadFileStreamRequest struct {
	Stream string `form:"stream"`
}

type UploadInitiateMultipartRequest struct {
	FileName string `form:"file_name" json:"file_name" binding:"required"`
	FileSize int64  `form:"file_size" json:"file_size" binding:"required"`
}

type UploadMultipartRequest struct {
	UploadId   string `form:"hash" json:"hash" binding:"required"`
	Name       string `form:"name" json:"name" binding:"required"`
	Size       int64  `form:"size" json:"size" binding:"required"`
	Ext        string `form:"ext" json:"ext" binding:"required"`
	SplitIndex int    `form:"split_index" json:"split_index" binding:"min=0"`
	SplitNum   int    `form:"split_num" json:"split_num" binding:"required,min=1"`
}
