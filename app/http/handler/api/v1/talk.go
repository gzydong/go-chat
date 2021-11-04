package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/entity"
	"go-chat/app/http/dto"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/im"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service"
	"strconv"
	"strings"
)

type Talk struct {
	service         *service.TalkService
	talkListService *service.TalkListService
	redisLock       *cache.RedisLock
	userService     *service.UserService
	wsClient        *cache.WsClient
	lastMessage     *cache.LastMessage
}

func NewTalkHandler(
	service *service.TalkService,
	talkListService *service.TalkListService,
	redisLock *cache.RedisLock,
	userService *service.UserService,
	wsClient *cache.WsClient,
	lastMessage *cache.LastMessage,
) *Talk {
	return &Talk{service, talkListService, redisLock, userService, wsClient, lastMessage}
}

// List 会话列表
func (c *Talk) List(ctx *gin.Context) {
	uid := auth.GetAuthUserID(ctx)

	data, err := c.talkListService.GetTalkList(ctx.Request.Context(), uid)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	items := make([]*dto.TalkListItem, 0)

	for _, item := range data {
		value := &dto.TalkListItem{
			ID:         item.ID,
			TalkType:   item.TalkType,
			ReceiverId: item.ReceiverId,
			IsTop:      item.IsTop,
			IsDisturb:  item.IsDisturb,
			IsRobot:    item.IsRobot,
			Avatar:     item.UserAvatar,
			MsgText:    "",
			UpdatedAt:  timeutil.FormatDatetime(item.UpdatedAt),
		}

		if item.TalkType == 1 {
			value.Name = item.Nickname
			value.Avatar = item.UserAvatar
			value.RemarkName = "" // 查询缓存
			value.UnreadNum = 0   // 查询缓存

			if c.wsClient.IsOnlineAll(ctx, im.GroupManage.DefaultChannel.Name, strconv.Itoa(value.ReceiverId)) {
				value.IsOnline = 1
			}
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		if msg, err := c.lastMessage.Get(ctx.Request.Context(), item.TalkType, uid, item.ReceiverId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		// 查询最后一条对话消息

		items = append(items, value)
	}

	response.Success(ctx, items)
}

// Create 创建会话列表
func (c *Talk) Create(ctx *gin.Context) {
	params := &request.TalkListCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	agent := strings.TrimSpace(ctx.GetHeader("user-agent"))
	if agent != "" {
		agent = strutil.Md5([]byte(agent))
	}

	lockKey := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.redisLock.Lock(ctx.Request.Context(), lockKey, 20) {
		response.BusinessError(ctx, "创建失败")
		return
	}

	result, err := c.talkListService.Create(ctx.Request.Context(), uid, params)
	if err != nil {
		response.BusinessError(ctx, err.Error())
		return
	}

	item := dto.TalkListItem{
		ID:         result.ID,
		TalkType:   result.TalkType,
		ReceiverId: result.ReceiverId,
		IsRobot:    result.IsRobot,
		Avatar:     "",
		Name:       "",
		RemarkName: "",
		UnreadNum:  0,
		MsgText:    "",
		UpdatedAt:  timeutil.DateTime(),
	}

	if item.TalkType == entity.PrivateChat {
		if user, err := c.userService.Dao().FindById(item.ReceiverId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else {

	}

	response.Success(ctx, &item)
}

// Delete 删除列表
func (c *Talk) Delete(ctx *gin.Context) {
	params := &request.TalkListDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Delete(ctx, auth.GetAuthUserID(ctx), params.Id); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Top 置顶列表
func (c *Talk) Top(ctx *gin.Context) {
	params := &request.TalkListTopRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Top(ctx, auth.GetAuthUserID(ctx), params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Disturb 会话免打扰
func (c *Talk) Disturb(ctx *gin.Context) {
	params := &request.TalkListDisturbRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Disturb(ctx, auth.GetAuthUserID(ctx), params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}
