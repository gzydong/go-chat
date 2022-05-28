package talk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/service"
)

type Talk struct {
	service         *service.TalkService
	talkListService *service.TalkSessionService
	redisLock       *cache.RedisLock
	userService     *service.UserService
	wsClient        *cache.WsClientSession
	lastMessage     *cache.LastMessage
	contactService  *service.ContactService
	unreadTalkCache *cache.UnreadTalkCache
	groupService    *service.GroupService
	authPermission  *service.AuthPermissionService
}

func NewTalkHandler(
	service *service.TalkService,
	talkListService *service.TalkSessionService,
	redisLock *cache.RedisLock,
	userService *service.UserService,
	wsClient *cache.WsClientSession,
	lastMessage *cache.LastMessage,
	unreadTalkCache *cache.UnreadTalkCache,
	contactService *service.ContactService,
	groupService *service.GroupService,
	authPermission *service.AuthPermissionService,
) *Talk {
	return &Talk{
		service:         service,
		talkListService: talkListService,
		redisLock:       redisLock,
		userService:     userService,
		wsClient:        wsClient,
		lastMessage:     lastMessage,
		unreadTalkCache: unreadTalkCache,
		contactService:  contactService,
		groupService:    groupService,
		authPermission:  authPermission,
	}
}

// List 会话列表
func (c *Talk) List(ctx *gin.Context) {
	uid := jwtutil.GetUid(ctx)

	data, err := c.talkListService.List(ctx.Request.Context(), uid)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkType == 1 {
			friends = append(friends, item.ReceiverId)
		}
	}

	remarks, err := c.contactService.Dao().GetFriendRemarks(ctx, uid, friends)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	items := make([]*dto.TalkListItem, 0)
	for _, item := range data {
		value := &dto.TalkListItem{
			Id:         item.Id,
			TalkType:   item.TalkType,
			ReceiverId: item.ReceiverId,
			IsTop:      item.IsTop,
			IsDisturb:  item.IsDisturb,
			IsRobot:    item.IsRobot,
			Avatar:     item.UserAvatar,
			MsgText:    "...",
			UpdatedAt:  timeutil.FormatDatetime(item.UpdatedAt),
		}

		// TODO 需要优化加缓存
		if item.TalkType == 1 {
			value.Name = item.Nickname
			value.Avatar = item.UserAvatar
			value.RemarkName = remarks[item.ReceiverId]
			value.UnreadNum = c.unreadTalkCache.Get(ctx.Request.Context(), item.ReceiverId, uid)
			value.IsOnline = strutil.BoolToInt(c.wsClient.IsOnline(ctx, entity.ImChannelDefault, strconv.Itoa(value.ReceiverId)))
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := c.lastMessage.Get(ctx.Request.Context(), item.TalkType, uid, item.ReceiverId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	response.Success(ctx, items)
}

// Create 创建会话列表
func (c *Talk) Create(ctx *gin.Context) {
	var (
		params = &request.TalkListCreateRequest{}
		uid    = jwtutil.GetUid(ctx)
		agent  = strings.TrimSpace(ctx.GetHeader("user-agent"))
	)

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if agent != "" {
		agent = encrypt.Md5(agent)
	}

	// 判断对方是否是自己
	if params.TalkType == entity.ChatPrivateMode && params.ReceiverId == jwtutil.GetUid(ctx) {
		response.BusinessError(ctx, "创建失败")
		return
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.redisLock.Lock(ctx.Request.Context(), key, 10) {
		response.BusinessError(ctx, "创建失败")
		return
	}

	// 暂无权限
	if !c.authPermission.IsAuth(ctx.Request.Context(), &service.AuthPermission{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "暂无权限！")
		return
	}

	result, err := c.talkListService.Create(ctx.Request.Context(), &service.TalkSessionCreateOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
	})
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	item := dto.TalkListItem{
		Id:         result.Id,
		TalkType:   result.TalkType,
		ReceiverId: result.ReceiverId,
		IsRobot:    result.IsRobot,
		UpdatedAt:  timeutil.DateTime(),
	}

	if item.TalkType == entity.ChatPrivateMode {
		item.UnreadNum = c.unreadTalkCache.Get(ctx.Request.Context(), params.ReceiverId, uid)
		item.RemarkName = c.contactService.Dao().GetFriendRemark(ctx.Request.Context(), uid, params.ReceiverId, true)

		if user, err := c.userService.Dao().FindById(item.ReceiverId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if item.TalkType == entity.ChatGroupMode {
		if group, err := c.groupService.Dao().FindById(params.ReceiverId); err == nil {
			item.Name = group.Name
		}
	}

	// 查询缓存消息
	if msg, err := c.lastMessage.Get(ctx.Request.Context(), item.TalkType, uid, item.ReceiverId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
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

	if err := c.talkListService.Delete(ctx, jwtutil.GetUid(ctx), params.Id); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Top 置顶列表
func (c *Talk) Top(ctx *gin.Context) {
	params := &request.TalkListTopRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Top(ctx, &service.TalkSessionTopOpts{
		UserId: jwtutil.GetUid(ctx),
		Id:     params.Id,
		Type:   params.Type,
	}); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Disturb 会话免打扰
func (c *Talk) Disturb(ctx *gin.Context) {
	params := &request.TalkListDisturbRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkListService.Disturb(ctx, &service.TalkSessionDisturbOpts{
		UserId:     jwtutil.GetUid(ctx),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		IsDisturb:  params.IsDisturb,
	}); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

func (c *Talk) ClearUnreadMessage(ctx *gin.Context) {
	params := &request.TalkUnReadRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)
	if params.TalkType == 1 {
		c.unreadTalkCache.Reset(ctx.Request.Context(), params.ReceiverId, uid)
	}

	response.Success(ctx, nil)
}
