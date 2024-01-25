package v1

import (
	"bytes"
	"slices"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

type Emoticon struct {
	RedisLock       *cache.RedisLock
	EmoticonRepo    *repo.Emoticon
	EmoticonService service.IEmoticonService
	Filesystem      filesystem.IFilesystem
}

// List 收藏列表
func (c *Emoticon) List(ctx *core.Context) error {
	resp := &web.EmoticonListResponse{
		Items: make([]*web.EmoticonItem, 0),
	}

	items, err := c.EmoticonRepo.GetCustomizeList(ctx.UserId())
	if err != nil {
		return ctx.Error(err.Error())
	}

	for _, item := range items {
		resp.Items = append(resp.Items, &web.EmoticonItem{
			EmoticonId: int32(item.Id),
			Url:        item.Url,
		})
	}

	return ctx.Success(resp)
}

// Delete 删除收藏表情包
func (c *Emoticon) Delete(ctx *core.Context) error {
	in := &web.EmoticonDeleteRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.EmoticonService.DeleteCollect(ctx.UserId(), []int{int(in.GetEmoticonId())}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// Create 创建自定义表情包
func (c *Emoticon) Create(ctx *core.Context) error {
	in := &web.EmoticonCreateRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	m := &model.EmoticonItem{
		UserId:   ctx.UserId(),
		Describe: "自定义表情包",
		Url:      in.Url,
	}

	if err := c.EmoticonRepo.Db.Create(m).Error; err != nil {
		return ctx.ErrorBusiness("上传失败！")
	}

	return ctx.Success(&web.EmoticonCreateResponse{
		EmoticonId: int32(m.Id),
		Url:        m.Url,
	})
}

// Upload 上传自定义表情包
func (c *Emoticon) Upload(ctx *core.Context) error {
	file, err := ctx.Context.FormFile("file")
	if err != nil {
		return ctx.InvalidParams("file 字段必传！")
	}

	if !slices.Contains([]string{"png", "jpg", "jpeg", "gif"}, strutil.FileSuffix(file.Filename)) {
		return ctx.InvalidParams("上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		return ctx.InvalidParams("上传文件大小不能超过5M！")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return ctx.ErrorBusiness("上传失败！")
	}

	meta := utils.ReadImageMeta(bytes.NewReader(stream))
	ext := strutil.FileSuffix(file.Filename)

	src := strutil.GenMediaObjectName(ext, meta.Width, meta.Height)
	if err = c.Filesystem.Write(c.Filesystem.BucketPublicName(), src, stream); err != nil {
		return ctx.ErrorBusiness("上传失败！")
	}

	m := &model.EmoticonItem{
		UserId:   ctx.UserId(),
		Describe: "自定义表情包",
		Url:      c.Filesystem.PublicUrl(c.Filesystem.BucketPublicName(), src),
	}

	if err := c.EmoticonRepo.Db.Create(m).Error; err != nil {
		return ctx.ErrorBusiness("上传失败！")
	}

	return ctx.Success(&web.EmoticonUploadResponse{
		EmoticonId: int32(m.Id),
		Url:        m.Url,
	})
}
