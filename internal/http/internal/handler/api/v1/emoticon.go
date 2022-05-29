package v1

import (
	"bytes"
	"fmt"
	"time"

	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/api"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/model"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/utils"

	"github.com/gin-gonic/gin"

	"go-chat/internal/http/internal/response"
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

func NewEmoticonHandler(
	service *service.EmoticonService,
	fileSystem *filesystem.Filesystem,
	redisLock *cache.RedisLock,
) *Emoticon {
	return &Emoticon{
		service:    service,
		fileSystem: fileSystem,
		redisLock:  redisLock,
	}
}

// CollectList 收藏列表
func (c *Emoticon) CollectList(ctx *gin.Context) {
	var (
		uid     = jwtutil.GetUid(ctx)
		sys     = make([]*api.SysEmoticonResponse, 0)
		collect = make([]*api.EmoticonItem, 0)
	)

	if ids := c.service.Dao().GetUserInstallIds(uid); len(ids) > 0 {
		if items, err := c.service.Dao().FindByIds(ids); err == nil {
			for _, item := range items {
				data := &api.SysEmoticonResponse{
					EmoticonId: item.Id,
					Url:        item.Icon,
					Name:       item.Name,
					List:       make([]*api.EmoticonItem, 0),
				}

				if items, err := c.service.Dao().GetDetailsAll(item.Id, 0); err == nil {
					for _, item := range items {
						data.List = append(data.List, &api.EmoticonItem{
							MediaId: item.Id,
							Src:     item.Url,
						})
					}
				}

				sys = append(sys, data)
			}
		}
	}

	if items, err := c.service.Dao().GetDetailsAll(0, uid); err == nil {
		for _, item := range items {
			collect = append(collect, &api.EmoticonItem{
				MediaId: item.Id,
				Src:     item.Url,
			})
		}
	}

	response.Success(ctx, entity.H{
		"sys_emoticon":     sys,
		"collect_emoticon": collect,
	})
}

// DeleteCollect 删除收藏表情包
func (c *Emoticon) DeleteCollect(ctx *gin.Context) {
	params := &request.DeleteCollectRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.DeleteCollect(jwtutil.GetUid(ctx), sliceutil.ParseIds(params.Ids)); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Upload 上传自定义表情包
func (c *Emoticon) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("emoticon")
	if err != nil {
		response.InvalidParams(ctx, "emoticon 字段必传！")
		return
	}

	if !sliceutil.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif"}) {
		response.InvalidParams(ctx, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(ctx, "上传文件大小不能超过5M！")
		return
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		response.BusinessError(ctx, "上传失败！")
		return
	}

	size := utils.ReadFileImage(bytes.NewReader(stream))
	ext := strutil.FileSuffix(file.Filename)
	src := fmt.Sprintf("public/media/image/emoticon/%s/%s", time.Now().Format("20060102"), strutil.GenImageName(ext, size.Width, size.Height))
	if err = c.fileSystem.Default.Write(stream, src); err != nil {
		response.BusinessError(ctx, "上传失败！")
		return
	}

	m := &model.EmoticonItem{
		UserId:     jwtutil.GetUid(ctx),
		Describe:   "自定义表情包",
		Url:        c.fileSystem.Default.PublicUrl(src),
		FileSuffix: ext,
		FileSize:   int(file.Size),
	}

	if err := c.service.Db().Create(m).Error; err != nil {
		response.BusinessError(ctx, "上传失败！")
		return
	}

	response.Success(ctx, entity.H{
		"media_id": m.Id,
		"src":      m.Url,
	}, "文件上传成功")
}

// SystemList 系统表情包列表
func (c *Emoticon) SystemList(ctx *gin.Context) {
	items, err := c.service.Dao().GetSystemEmoticonList()

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	ids := c.service.Dao().GetUserInstallIds(jwtutil.GetUid(ctx))

	data := make([]*api.SysEmoticonList, 0, len(items))
	for _, item := range items {
		data = append(data, &api.SysEmoticonList{
			ID:     item.Id,
			Name:   item.Name,
			Icon:   item.Icon,
			Status: strutil.BoolToInt(sliceutil.InInt(item.Id, ids)), // 查询用户是否使用
		})
	}

	response.Success(ctx, data)
}

// SetSystemEmoticon 添加或移除系统表情包
func (c *Emoticon) SetSystemEmoticon(ctx *gin.Context) {
	var (
		err    error
		params = &request.SetSystemEmoticonRequest{}
		uid    = jwtutil.GetUid(ctx)
		key    = fmt.Sprintf("sys-emoticon:%d", uid)
	)

	if err = ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !c.redisLock.Lock(ctx, key, 5) {
		response.BusinessError(ctx, "请求频繁！")
		return
	}

	defer c.redisLock.UnLock(ctx, key)

	if params.Type == 2 {
		if err = c.service.RemoveUserSysEmoticon(uid, params.EmoticonId); err != nil {
			response.BusinessError(ctx, err)
		} else {
			response.Success(ctx, nil)
		}

		return
	}

	// 查询表情包是否存在
	info, err := c.service.Dao().FindById(params.EmoticonId)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	if err = c.service.AddUserSysEmoticon(uid, params.EmoticonId); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	items := make([]*api.EmoticonItem, 0)
	if list, err := c.service.Dao().GetDetailsAll(params.EmoticonId, 0); err == nil {
		for _, item := range list {
			items = append(items, &api.EmoticonItem{
				MediaId: item.Id,
				Src:     item.Url,
			})
		}
	}

	response.Success(ctx, &api.SysEmoticonResponse{
		EmoticonId: info.Id,
		Url:        info.Icon,
		Name:       info.Name,
		List:       items,
	})
}
