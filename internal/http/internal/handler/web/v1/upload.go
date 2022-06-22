package v1

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ginutil"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
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
func (u *Upload) Avatar(ctx *gin.Context) error {

	file, err := ctx.FormFile("file")
	if err != nil {
		return ginutil.InvalidParams(ctx, "文件上传失败！")
	}

	stream, _ := filesystem.ReadMultipartStream(file)
	object := fmt.Sprintf("public/media/image/avatar/%s/%s", time.Now().Format("20060102"), strutil.GenImageName("png", 200, 200))

	if err := u.filesystem.Default.Write(stream, object); err != nil {
		return ginutil.BusinessError(ctx, "文件上传失败")
	}

	return ginutil.Success(ctx, entity.H{
		"avatar": u.filesystem.Default.PublicUrl(object),
	})
}

// InitiateMultipart 批量上传初始化
func (u *Upload) InitiateMultipart(ctx *gin.Context) error {
	params := &web.UploadInitiateMultipartRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	info, err := u.service.InitiateMultipartUpload(ctx.Request.Context(), &service.MultipartInitiateOpts{
		Name:   params.FileName,
		Size:   params.FileSize,
		UserId: jwtutil.GetUid(ctx),
	})
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, entity.H{
		"upload_id":  info.UploadId,
		"split_size": 2 << 20,
	})
}

// MultipartUpload 批量分片上传
func (u *Upload) MultipartUpload(ctx *gin.Context) error {
	params := &web.UploadMultipartRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return ginutil.InvalidParams(ctx, "文件上传失败！")
	}

	err = u.service.MultipartUpload(ctx.Request.Context(), &service.MultipartUploadOpts{
		UserId:     jwtutil.GetUid(ctx),
		UploadId:   params.UploadId,
		SplitIndex: params.SplitIndex,
		SplitNum:   params.SplitNum,
		File:       file,
	})
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	if params.SplitIndex != params.SplitNum-1 {
		return ginutil.Success(ctx, entity.H{"is_merge": false})
	}

	return ginutil.Success(ctx, entity.H{"is_merge": true, "upload_id": params.UploadId})
}
