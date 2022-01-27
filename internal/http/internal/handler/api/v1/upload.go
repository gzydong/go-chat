package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
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

// Avatar 头像上传上传
func (u *Upload) Avatar(ctx *gin.Context) {

	file, err := ctx.FormFile("file")
	if err != nil {
		response.InvalidParams(ctx, "文件上传失败！")
		return
	}

	stream, _ := filesystem.ReadMultipartStream(file)
	object := fmt.Sprintf("public/media/image/avatar/%s/%s", time.Now().Format("20060102"), strutil.GenImageName("png", 200, 200))

	if err := u.filesystem.Default.Write(stream, object); err != nil {
		response.BusinessError(ctx, "文件上传失败")
		return
	}

	response.Success(ctx, gin.H{
		"avatar": u.filesystem.Default.PublicUrl(object),
	})
}

// InitiateMultipart 批量上传初始化
func (u *Upload) InitiateMultipart(ctx *gin.Context) {
	params := &request.UploadInitiateMultipartRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	info, err := u.service.InitiateMultipartUpload(ctx.Request.Context(), &service.MultipartInitiateOpts{
		Name:   params.FileName,
		Size:   params.FileSize,
		UserId: jwt.GetUid(ctx),
	})
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, &gin.H{
		"upload_id":  info.UploadId,
		"split_size": 2 << 20,
	})
}

// MultipartUpload 批量分片上传
func (u *Upload) MultipartUpload(ctx *gin.Context) {
	params := &request.UploadMultipartRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.InvalidParams(ctx, "文件上传失败！")
		return
	}

	err = u.service.MultipartUpload(ctx.Request.Context(), &service.MultipartUploadOpts{
		UserId:     jwt.GetUid(ctx),
		UploadId:   params.UploadId,
		SplitIndex: params.SplitIndex,
		SplitNum:   params.SplitNum,
		File:       file,
	})
	if err != nil {
		response.BusinessError(ctx, err)
	}

	if params.SplitIndex != params.SplitNum-1 {
		response.Success(ctx, gin.H{"is_merge": false})
	} else {
		response.Success(ctx, gin.H{"is_merge": true, "upload_id": params.UploadId})
	}
}
