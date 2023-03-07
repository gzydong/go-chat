package v1

import (
	"bytes"
	"fmt"
	"time"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

type Emoticon struct {
	fileSystem *filesystem.Filesystem
	service    *service.EmoticonService
	redisLock  *cache.RedisLock
}

func NewEmoticon(fileSystem *filesystem.Filesystem, service *service.EmoticonService, redisLock *cache.RedisLock) *Emoticon {
	return &Emoticon{fileSystem: fileSystem, service: service, redisLock: redisLock}
}

// CollectList 收藏列表
func (c *Emoticon) CollectList(ctx *ichat.Context) error {

	var (
		uid  = ctx.UserId()
		resp = &web.EmoticonListResponse{
			SysEmoticon:     make([]*web.EmoticonListResponse_SysEmoticon, 0),
			CollectEmoticon: make([]*web.EmoticonListItem, 0),
		}
	)

	if ids := c.service.Dao().GetUserInstallIds(uid); len(ids) > 0 {
		if items, err := c.service.Dao().FindByIds(ctx.Ctx(), ids); err == nil {
			for _, item := range items {
				data := &web.EmoticonListResponse_SysEmoticon{
					EmoticonId: int32(item.Id),
					Url:        item.Icon,
					Name:       item.Name,
					List:       make([]*web.EmoticonListItem, 0),
				}

				if list, err := c.service.Dao().GetDetailsAll(item.Id, 0); err == nil {
					for _, v := range list {
						data.List = append(data.List, &web.EmoticonListItem{
							MediaId: int32(v.Id),
							Src:     v.Url,
						})
					}
				}

				resp.SysEmoticon = append(resp.SysEmoticon, data)
			}
		}
	}

	if items, err := c.service.Dao().GetDetailsAll(0, uid); err == nil {
		for _, item := range items {
			resp.CollectEmoticon = append(resp.CollectEmoticon, &web.EmoticonListItem{
				MediaId: int32(item.Id),
				Src:     item.Url,
			})
		}
	}

	return ctx.Success(resp)
}

// DeleteCollect 删除收藏表情包
func (c *Emoticon) DeleteCollect(ctx *ichat.Context) error {

	params := &web.EmoticonDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.DeleteCollect(ctx.UserId(), sliceutil.ParseIds(params.Ids)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// Upload 上传自定义表情包
func (c *Emoticon) Upload(ctx *ichat.Context) error {

	file, err := ctx.Context.FormFile("emoticon")
	if err != nil {
		return ctx.InvalidParams("emoticon 字段必传！")
	}

	if !sliceutil.Include(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif"}) {
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
	src := fmt.Sprintf("public/media/image/emoticon/%s/%s", time.Now().Format("20060102"), strutil.GenImageName(ext, meta.Width, meta.Height))
	if err = c.fileSystem.Default.Write(stream, src); err != nil {
		return ctx.ErrorBusiness("上传失败！")
	}

	m := &model.EmoticonItem{
		UserId:     ctx.UserId(),
		Describe:   "自定义表情包",
		Url:        c.fileSystem.Default.PublicUrl(src),
		FileSuffix: ext,
		FileSize:   int(file.Size),
	}

	if err := c.service.Db().Create(m).Error; err != nil {
		return ctx.ErrorBusiness("上传失败！")
	}

	return ctx.Success(&web.EmoticonUploadResponse{
		MediaId: int32(m.Id),
		Src:     m.Url,
	})
}

// SystemList 系统表情包列表
func (c *Emoticon) SystemList(ctx *ichat.Context) error {

	items, err := c.service.Dao().GetSystemEmoticonList(ctx.Ctx())

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	ids := c.service.Dao().GetUserInstallIds(ctx.UserId())

	data := make([]*web.EmoticonSysListResponse_Item, 0)
	for _, item := range items {
		data = append(data, &web.EmoticonSysListResponse_Item{
			Id:     int32(item.Id),
			Name:   item.Name,
			Icon:   item.Icon,
			Status: int32(strutil.BoolToInt(sliceutil.Include(item.Id, ids))), // 查询用户是否使用
		})
	}

	return ctx.Success(data)
}

// SetSystemEmoticon 添加或移除系统表情包
func (c *Emoticon) SetSystemEmoticon(ctx *ichat.Context) error {
	var (
		err    error
		params = &web.EmoticonSetSystemRequest{}
		uid    = ctx.UserId()
		key    = fmt.Sprintf("sys-emoticon:%d", uid)
	)

	if err = ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if !c.redisLock.Lock(ctx.Ctx(), key, 5) {
		return ctx.ErrorBusiness("请求频繁！")
	}
	defer c.redisLock.UnLock(ctx.Ctx(), key)

	if params.Type == 2 {
		if err = c.service.RemoveUserSysEmoticon(uid, int(params.EmoticonId)); err != nil {
			return ctx.ErrorBusiness(err.Error())
		}

		return ctx.Success(nil)
	}

	// 查询表情包是否存在
	info, err := c.service.Dao().FindById(ctx.Ctx(), int(params.EmoticonId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.service.AddUserSysEmoticon(uid, int(params.EmoticonId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	items := make([]*web.EmoticonListItem, 0)
	if list, err := c.service.Dao().GetDetailsAll(int(params.EmoticonId), 0); err == nil {
		for _, item := range list {
			items = append(items, &web.EmoticonListItem{
				MediaId: int32(item.Id),
				Src:     item.Url,
			})
		}
	}

	return ctx.Success(&web.EmoticonSetSystemResponse{
		EmoticonId: int32(info.Id),
		Url:        info.Icon,
		Name:       info.Name,
		List:       items,
	})
}
