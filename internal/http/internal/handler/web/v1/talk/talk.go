package talk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/ginutil"
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
func (c *Talk) List(ctx *gin.Context) error {
	uid := jwtutil.GetUid(ctx)

	data, err := c.talkListService.List(ctx.Request.Context(), uid)
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkType == 1 {
			friends = append(friends, item.ReceiverId)
		}
	}

	remarks, err := c.contactService.Dao().GetFriendRemarks(ctx, uid, friends)
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	items := make([]*web.TalkListItem, 0)
	for _, item := range data {
		value := &web.TalkListItem{
			Id:         int32(item.Id),
			TalkType:   int32(item.TalkType),
			ReceiverId: int32(item.ReceiverId),
			IsTop:      int32(item.IsTop),
			IsDisturb:  int32(item.IsDisturb),
			IsRobot:    int32(item.IsRobot),
			Avatar:     item.UserAvatar,
			MsgText:    "...",
			UpdatedAt:  timeutil.FormatDatetime(item.UpdatedAt),
		}

		// TODO 需要优化加缓存
		if item.TalkType == 1 {
			value.Name = item.Nickname
			value.Avatar = item.UserAvatar
			value.RemarkName = remarks[item.ReceiverId]
			value.UnreadNum = int32(c.unreadTalkCache.Get(ctx.Request.Context(), item.ReceiverId, uid))
			value.IsOnline = int32(strutil.BoolToInt(c.wsClient.IsOnline(ctx, entity.ImChannelDefault, strconv.Itoa(int(value.ReceiverId)))))
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

	return ginutil.Success(ctx, &web.GetTalkListResponse{
		Items: items,
	})
}

// Create 创建会话列表
func (c *Talk) Create(ctx *gin.Context) error {
	var (
		params = &web.CreateTalkListRequest{}
		uid    = jwtutil.GetUid(ctx)
		agent  = strings.TrimSpace(ctx.GetHeader("user-agent"))
	)

	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	if agent != "" {
		agent = encrypt.Md5(agent)
	}

	// 判断对方是否是自己
	if params.TalkType == entity.ChatPrivateMode && params.ReceiverId == jwtutil.GetUid(ctx) {
		return ginutil.BusinessError(ctx, "创建失败")
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.redisLock.Lock(ctx.Request.Context(), key, 10) {
		return ginutil.BusinessError(ctx, "创建失败")
	}

	// 暂无权限
	if !c.authPermission.IsAuth(ctx.Request.Context(), &service.AuthPermission{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		return ginutil.BusinessError(ctx, "暂无权限！")
	}

	result, err := c.talkListService.Create(ctx.Request.Context(), &service.TalkSessionCreateOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
	})
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	item := &web.TalkListItem{
		Id:         int32(result.Id),
		TalkType:   int32(result.TalkType),
		ReceiverId: int32(result.ReceiverId),
		IsRobot:    int32(result.IsRobot),
		UpdatedAt:  timeutil.DateTime(),
	}

	if item.TalkType == entity.ChatPrivateMode {
		item.UnreadNum = int32(c.unreadTalkCache.Get(ctx.Request.Context(), params.ReceiverId, uid))
		item.RemarkName = c.contactService.Dao().GetFriendRemark(ctx.Request.Context(), uid, params.ReceiverId, true)

		if user, err := c.userService.Dao().FindById(result.ReceiverId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if result.TalkType == entity.ChatGroupMode {
		if group, err := c.groupService.Dao().FindById(params.ReceiverId); err == nil {
			item.Name = group.Name
		}
	}

	// 查询缓存消息
	if msg, err := c.lastMessage.Get(ctx.Request.Context(), result.TalkType, uid, result.ReceiverId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return ginutil.Success(ctx, &web.CreateTalkListResponse{
		Item: item,
	})
}

// Delete 删除列表
func (c *Talk) Delete(ctx *gin.Context) error {
	params := &web.DeleteTalkListRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	if err := c.talkListService.Delete(ctx, jwtutil.GetUid(ctx), params.Id); err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, &web.DeleteTalkListResponse{})
}

// Top 置顶列表
func (c *Talk) Top(ctx *gin.Context) error {
	params := &web.TopTalkListRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	if err := c.talkListService.Top(ctx, &service.TalkSessionTopOpts{
		UserId: jwtutil.GetUid(ctx),
		Id:     params.Id,
		Type:   params.Type,
	}); err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, &web.TopTalkListResponse{})
}

// Disturb 会话免打扰
func (c *Talk) Disturb(ctx *gin.Context) error {
	params := &web.DisturbTalkListRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	if err := c.talkListService.Disturb(ctx, &service.TalkSessionDisturbOpts{
		UserId:     jwtutil.GetUid(ctx),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		IsDisturb:  params.IsDisturb,
	}); err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, &web.DisturbTalkListResponse{})
}

func (c *Talk) ClearUnreadMessage(ctx *gin.Context) error {
	params := &web.ClearTalkUnreadNumRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)
	if params.TalkType == 1 {
		c.unreadTalkCache.Reset(ctx.Request.Context(), params.ReceiverId, uid)
	}

	return ginutil.Success(ctx, &web.ClearTalkUnreadNumResponse{})
}
