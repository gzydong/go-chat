package v1

import (
	"bytes"
	"math"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-chat/api/pb/web/v1"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
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

// Image 图片上传
func (u *Upload) Image(ctx *gin.Context) (*web.UploadImageResponse, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}

	var (
		ext       = strings.TrimPrefix(path.Ext(file.Filename), ".")
		width, _  = strconv.Atoi(ctx.DefaultPostForm("width", "0"))
		height, _ = strconv.Atoi(ctx.DefaultPostForm("height", "0"))
	)

	stream, _ := filesystem.ReadMultipartStream(file)
	if width == 0 || height == 0 {
		meta := utils.ReadImageMeta(bytes.NewReader(stream))
		width = meta.Width
		height = meta.Height
	}

	object := strutil.GenMediaObjectName(ext, width, height)
	if err := u.Filesystem.Write(u.Filesystem.BucketPublicName(), object, stream); err != nil {
		return nil, err
	}

	return &web.UploadImageResponse{
		Src: u.Filesystem.PublicUrl(u.Filesystem.BucketPublicName(), object),
	}, nil
}

// InitiateMultipart 批量上传初始化
func (u *Upload) InitiateMultipart(ctx *gin.Context) (*web.UploadInitiateMultipartResponse, error) {
	in := &web.UploadInitiateMultipartRequest{}
	if err := ctx.ShouldBindJSON(in); err != nil {
		return nil, errorx.New(400, err.Error())
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx.Request.Context())
	info, err := u.SplitUploadService.InitiateMultipartUpload(ctx, &service.MultipartInitiateOpt{
		Name:   in.FileName,
		Size:   in.FileSize,
		UserId: uid,
	})

	if err != nil {
		return nil, err
	}

	return &web.UploadInitiateMultipartResponse{
		UploadId:  info.UploadId,
		ShardSize: 5 << 20,
		ShardNum:  int32(math.Ceil(float64(in.FileSize) / float64(5<<20))),
	}, nil
}

// MultipartUpload 批量分片上传
func (u *Upload) MultipartUpload(ctx *gin.Context) (*web.UploadMultipartResponse, error) {
	in := &web.UploadMultipartRequest{
		UploadId: ctx.PostForm("upload_id"),
	}

	splitIndex, err := strconv.Atoi(ctx.PostForm("split_index"))
	if err != nil {
		return nil, errorx.New(400, "split_index 不合法")
	}

	splitNum, err := strconv.Atoi(ctx.PostForm("split_num"))
	if err != nil {
		return nil, errorx.New(400, "split_num 不合法")
	}

	in.SplitIndex = int32(splitIndex)
	in.SplitNum = int32(splitNum)

	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, errorx.New(400, "文件上传失败")
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx.Request.Context())
	err = u.SplitUploadService.MultipartUpload(ctx.Request.Context(), &service.MultipartUploadOpt{
		UserId:     uid,
		UploadId:   in.UploadId,
		SplitIndex: int(in.SplitIndex),
		SplitNum:   int(in.SplitNum),
		File:       file,
	})
	if err != nil {
		return nil, err
	}

	if in.SplitIndex != in.SplitNum {
		return &web.UploadMultipartResponse{
			IsMerge: false,
		}, nil
	}

	return &web.UploadMultipartResponse{
		UploadId: in.UploadId,
		IsMerge:  true,
	}, nil
}
