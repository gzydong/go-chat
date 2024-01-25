package v1

import (
	"bytes"
	"math"
	"path"
	"strconv"
	"strings"

	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/service"
)

type Upload struct {
	Config             *config.Config
	Filesystem         filesystem.IFilesystem
	SplitUploadService service.ISplitUploadService
}

// Avatar 头像上传上传
func (u *Upload) Avatar(ctx *core.Context) error {
	file, err := ctx.Context.FormFile("file")
	if err != nil {
		return ctx.InvalidParams("文件上传失败！")
	}

	stream, _ := filesystem.ReadMultipartStream(file)

	object := strutil.GenMediaObjectName("png", 200, 200)
	if err := u.Filesystem.Write(u.Filesystem.BucketPublicName(), object, stream); err != nil {
		return ctx.ErrorBusiness("文件上传失败")
	}

	return ctx.Success(web.UploadAvatarResponse{
		Avatar: u.Filesystem.PublicUrl(u.Filesystem.BucketPublicName(), object),
	})
}

// Image 图片上传
func (u *Upload) Image(ctx *core.Context) error {

	file, err := ctx.Context.FormFile("file")
	if err != nil {
		return ctx.InvalidParams("文件上传失败！")
	}

	var (
		ext       = strings.TrimPrefix(path.Ext(file.Filename), ".")
		width, _  = strconv.Atoi(ctx.Context.DefaultPostForm("width", "0"))
		height, _ = strconv.Atoi(ctx.Context.DefaultPostForm("height", "0"))
	)

	stream, _ := filesystem.ReadMultipartStream(file)
	if width == 0 || height == 0 {
		meta := utils.ReadImageMeta(bytes.NewReader(stream))
		width = meta.Width
		height = meta.Height
	}

	object := strutil.GenMediaObjectName(ext, width, height)
	if err := u.Filesystem.Write(u.Filesystem.BucketPublicName(), object, stream); err != nil {
		return ctx.ErrorBusiness("文件上传失败")
	}

	return ctx.Success(web.UploadImageResponse{
		Src: u.Filesystem.PublicUrl(u.Filesystem.BucketPublicName(), object),
	})
}

// InitiateMultipart 批量上传初始化
func (u *Upload) InitiateMultipart(ctx *core.Context) error {
	in := &web.UploadInitiateMultipartRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	info, err := u.SplitUploadService.InitiateMultipartUpload(ctx.Ctx(), &service.MultipartInitiateOpt{
		Name:   in.FileName,
		Size:   in.FileSize,
		UserId: ctx.UserId(),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.UploadInitiateMultipartResponse{
		UploadId:  info.UploadId,
		ShardSize: 5 << 20,
		ShardNum:  int32(math.Ceil(float64(in.FileSize) / float64(5<<20))),
	})
}

// MultipartUpload 批量分片上传
func (u *Upload) MultipartUpload(ctx *core.Context) error {
	in := &web.UploadMultipartRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	file, err := ctx.Context.FormFile("file")
	if err != nil {
		return ctx.InvalidParams("文件上传失败！")
	}

	err = u.SplitUploadService.MultipartUpload(ctx.Ctx(), &service.MultipartUploadOpt{
		UserId:     ctx.UserId(),
		UploadId:   in.UploadId,
		SplitIndex: int(in.SplitIndex),
		SplitNum:   int(in.SplitNum),
		File:       file,
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if in.SplitIndex != in.SplitNum {
		return ctx.Success(&web.UploadMultipartResponse{
			IsMerge: false,
		})
	}

	return ctx.Success(&web.UploadMultipartResponse{
		UploadId: in.UploadId,
		IsMerge:  true,
	})
}
