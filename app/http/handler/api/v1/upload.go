package v1

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/strutil"
	"go-chat/app/service"
	"go-chat/config"
	"strings"
	"time"
)

type Upload struct {
	config     *config.Config
	filesystem *filesystem.Filesystem
	service    *service.SplitUploadService
}

func NewUploadHandler(
	config *config.Config,
	filesystem *filesystem.Filesystem,
	service *service.SplitUploadService,
) *Upload {
	return &Upload{
		config:     config,
		filesystem: filesystem,
		service:    service,
	}
}

// 文件流上传
func (u *Upload) Stream(ctx *gin.Context) {
	params := &request.UploadFileStreamRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	params.Stream = strings.Replace(params.Stream, "data:image/png;base64,", "", 1)
	params.Stream = strings.Replace(params.Stream, " ", "+", 1)

	stream, _ := base64.StdEncoding.DecodeString(params.Stream)

	object := fmt.Sprintf("public/media/image/avatar/%s/%s", time.Now().Format("20060102"), strutil.GenImageName("png", 200, 200))

	err := u.filesystem.Write(stream, object)
	if err != nil {
		response.BusinessError(ctx, "文件上传失败")
		return
	}

	response.Success(ctx, gin.H{
		"avatar": u.filesystem.PublicUrl(object),
	})
}

// 批量上传初始化
func (u *Upload) InitiateMultipart(ctx *gin.Context) {
	params := &request.UploadInitiateMultipartRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	info, err := u.service.InitiateMultipartUpload(ctx.Request.Context(), &service.InitiateParams{
		Name:   params.FileName,
		Size:   params.FileSize,
		UserId: auth.GetAuthUserID(ctx),
	})
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, &gin.H{
		"file_type":     info.FileType,
		"user_id":       info.UserId,
		"original_name": info.OriginalName,
		"hash_name":     info.HashName,
		"file_ext":      info.FileExt,
		"file_size":     info.FileSize,
		"split_num":     info.SplitNum,
		"split_index":   info.SplitIndex,
	})
}

// 批量分片上传
func (u *Upload) MultipartUpload(ctx *gin.Context) {

}
