package v1

import (
	"bytes"
	"path"
	"strconv"
	"strings"

	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
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
func (u *Upload) Avatar(ctx *ichat.Context) error {
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
func (u *Upload) Image(ctx *ichat.Context) error {

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
func (u *Upload) InitiateMultipart(ctx *ichat.Context) error {

	params := &web.UploadInitiateMultipartRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	info, err := u.SplitUploadService.InitiateMultipartUpload(ctx.Ctx(), &service.MultipartInitiateOpt{
		Name:   params.FileName,
		Size:   params.FileSize,
		UserId: ctx.UserId(),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.UploadInitiateMultipartResponse{
		UploadId:    info.UploadId,
		UploadIdMd5: encrypt.Md5(info.UploadId),
		SplitSize:   5 << 20,
	})
}

// MultipartUpload 批量分片上传
func (u *Upload) MultipartUpload(ctx *ichat.Context) error {

	params := &web.UploadMultipartRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	file, err := ctx.Context.FormFile("file")
	if err != nil {
		return ctx.InvalidParams("文件上传失败！")
	}

	err = u.SplitUploadService.MultipartUpload(ctx.Ctx(), &service.MultipartUploadOpt{
		UserId:     ctx.UserId(),
		UploadId:   params.UploadId,
		SplitIndex: int(params.SplitIndex),
		SplitNum:   int(params.SplitNum),
		File:       file,
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if params.SplitIndex != params.SplitNum {
		return ctx.Success(&web.UploadMultipartResponse{
			IsMerge: false,
		})
	}

	return ctx.Success(&web.UploadMultipartResponse{
		UploadId: params.UploadId,
		IsMerge:  true,
	})
}
