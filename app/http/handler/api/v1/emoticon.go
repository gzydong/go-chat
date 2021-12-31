package v1

import (
	"fmt"
	"go-chat/app/cache"
	"go-chat/app/http/dto/api"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"

	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/service"
)

type Emoticon struct {
	filesystem *filesystem.Filesystem
	service    *service.EmoticonService
	redisLock  *cache.RedisLock
}

func NewEmoticonHandler(
	service *service.EmoticonService,
	filesystem *filesystem.Filesystem,
	redisLock *cache.RedisLock,
) *Emoticon {
	return &Emoticon{
		service:    service,
		filesystem: filesystem,
		redisLock:  redisLock,
	}
}

// CollectList 收藏列表
func (c *Emoticon) CollectList(ctx *gin.Context) {
	var (
		uid     = auth.GetAuthUserID(ctx)
		sys     = make([]*api.SysEmoticonResponse, 0)
		collect = make([]*api.EmoticonItem, 0)
	)

	if ids := c.service.Dao().GetUserInstallIds(uid); len(ids) > 0 {
		items := make([]*model.Emoticon, 0)

		if _, err := c.service.Dao().FindByIds(&items, ids, "*"); err == nil {
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

	response.Success(ctx, gin.H{
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

	if err := c.service.DeleteCollect(auth.GetAuthUserID(ctx), slice.ParseIds(params.Ids)); err != nil {
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

	if !slice.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif"}) {
		response.InvalidParams(ctx, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(ctx, "上传文件大小不能超过5M！")
		return
	}

	info, err := c.service.CustomizeUpload(ctx.Request.Context(), auth.GetAuthUserID(ctx), file)
	if err != nil {
		response.BusinessError(ctx, "文件上传失败！")
		return
	}

	response.Success(ctx, gin.H{
		"media_id": info.Id,
		"src":      info.Url,
	}, "文件上传成功")
}

// SystemList 系统表情包列表
func (c *Emoticon) SystemList(ctx *gin.Context) {
	items, err := c.service.Dao().GetSystemEmoticonList()

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	ids := c.service.Dao().GetUserInstallIds(auth.GetAuthUserID(ctx))

	data := make([]*api.SysEmoticonList, 0, len(items))
	for _, item := range items {
		data = append(data, &api.SysEmoticonList{
			ID:     item.Id,
			Name:   item.Name,
			Icon:   item.Icon,
			Status: strutil.BoolToInt(slice.InInt(item.Id, ids)), // 查询用户是否使用
		})
	}

	response.Success(ctx, data)
}

// SetSystemEmoticon 添加或移除系统表情包
func (c *Emoticon) SetSystemEmoticon(ctx *gin.Context) {
	var (
		err    error
		params = &request.SetSystemEmoticonRequest{}
		uid    = auth.GetAuthUserID(ctx)
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

	defer c.redisLock.Release(ctx, key)

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
