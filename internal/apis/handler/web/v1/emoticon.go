package v1

import (
	"bytes"
	"context"
	"slices"

	"github.com/gin-gonic/gin"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

var _ web.IEmoticonHandler = (*Emoticon)(nil)

type Emoticon struct {
	RedisLock       *cache.RedisLock
	EmoticonRepo    *repo.Emoticon
	EmoticonService service.IEmoticonService
	Filesystem      filesystem.IFilesystem
}

// List 收藏列表
func (c *Emoticon) List(ctx context.Context, req *web.EmoticonListRequest) (*web.EmoticonListResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	resp := &web.EmoticonListResponse{
		Items: make([]*web.EmoticonItem, 0),
	}

	items, err := c.EmoticonRepo.GetCustomizeList(session.GetAuthID())
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		resp.Items = append(resp.Items, &web.EmoticonItem{
			EmoticonId: int32(item.Id),
			Url:        item.Url,
		})
	}

	return resp, nil
}

// Delete 删除收藏表情包
func (c *Emoticon) Delete(ctx context.Context, in *web.EmoticonDeleteRequest) (*web.EmoticonDeleteResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)
	return nil, c.EmoticonService.DeleteCollect(session.GetAuthID(), []int{int(in.GetEmoticonId())})
}

// Create 创建自定义表情包
func (c *Emoticon) Create(ctx context.Context, in *web.EmoticonCreateRequest) (*web.EmoticonCreateResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	m := &model.EmoticonItem{
		UserId:   session.GetAuthID(),
		Describe: "自定义表情包",
		Url:      in.Url,
	}

	if err := c.EmoticonRepo.Db.Create(m).Error; err != nil {
		return nil, err
	}

	return &web.EmoticonCreateResponse{
		EmoticonId: int32(m.Id),
		Url:        m.Url,
	}, nil
}

// Upload 上传自定义表情包
func (c *Emoticon) Upload(ginContext *gin.Context, _ *web.EmoticonUploadRequest) (*web.EmoticonUploadResponse, error) {
	ctx := ginContext.Request.Context()
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	file, err := ginContext.FormFile("file")
	if err != nil {
		return nil, errorx.New(400, "file 字段必传")
	}

	if !slices.Contains([]string{"png", "jpg", "jpeg", "gif"}, strutil.FileSuffix(file.Filename)) {
		return nil, errorx.New(400, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		return nil, errorx.New(400, "上传文件大小不能超过5M")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return nil, err
	}

	meta := utils.ReadImageMeta(bytes.NewReader(stream))
	ext := strutil.FileSuffix(file.Filename)

	src := strutil.GenMediaObjectName(ext, meta.Width, meta.Height)
	if err = c.Filesystem.Write(c.Filesystem.BucketPublicName(), src, stream); err != nil {
		return nil, err
	}

	m := &model.EmoticonItem{
		UserId:   session.GetAuthID(),
		Describe: "自定义表情包",
		Url:      c.Filesystem.PublicUrl(c.Filesystem.BucketPublicName(), src),
	}

	if err := c.EmoticonRepo.Db.Create(m).Error; err != nil {
		return nil, err
	}

	return &web.EmoticonUploadResponse{
		EmoticonId: int32(m.Id),
		Url:        m.Url,
	}, nil
}
