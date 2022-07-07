package talk

import (
	"fmt"
	"strconv"
	"strings"

	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
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

func NewTalk(service *service.TalkService, talkListService *service.TalkSessionService, redisLock *cache.RedisLock, userService *service.UserService, wsClient *cache.WsClientSession, lastMessage *cache.LastMessage, contactService *service.ContactService, unreadTalkCache *cache.UnreadTalkCache, groupService *service.GroupService, authPermission *service.AuthPermissionService) *Talk {
	return &Talk{service: service, talkListService: talkListService, redisLock: redisLock, userService: userService, wsClient: wsClient, lastMessage: lastMessage, contactService: contactService, unreadTalkCache: unreadTalkCache, groupService: groupService, authPermission: authPermission}
}

// List 会话列表
func (c *Talk) List(ctx *ichat.Context) error {

	uid := ctx.UserId()

	data, err := c.talkListService.List(ctx.RequestContext(), uid)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkType == 1 {
			friends = append(friends, item.ReceiverId)
		}
	}

	remarks, err := c.contactService.Dao().GetFriendRemarks(ctx.Context, uid, friends)
	if err != nil {
		return ctx.BusinessError(err.Error())
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
			value.UnreadNum = int32(c.unreadTalkCache.Get(ctx.RequestContext(), item.ReceiverId, uid))
			value.IsOnline = int32(strutil.BoolToInt(c.wsClient.IsOnline(ctx.Context, entity.ImChannelDefault, strconv.Itoa(int(value.ReceiverId)))))
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := c.lastMessage.Get(ctx.RequestContext(), item.TalkType, uid, item.ReceiverId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	return ctx.Success(&web.GetTalkListResponse{
		Items: items,
	})
}

// Create 创建会话列表
func (c *Talk) Create(ctx *ichat.Context) error {

	var (
		params = &web.CreateTalkListRequest{}
		uid    = ctx.UserId()
		agent  = strings.TrimSpace(ctx.Context.GetHeader("user-agent"))
	)

	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if agent != "" {
		agent = encrypt.Md5(agent)
	}

	// 判断对方是否是自己
	if params.TalkType == entity.ChatPrivateMode && params.ReceiverId == ctx.UserId() {
		return ctx.BusinessError("创建失败")
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.redisLock.Lock(ctx.RequestContext(), key, 10) {
		return ctx.BusinessError("创建失败")
	}

	// 暂无权限
	if !c.authPermission.IsAuth(ctx.RequestContext(), &service.AuthPermission{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		return ctx.BusinessError("暂无权限！")
	}

	result, err := c.talkListService.Create(ctx.RequestContext(), &service.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
	})
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	item := &web.TalkListItem{
		Id:         int32(result.Id),
		TalkType:   int32(result.TalkType),
		ReceiverId: int32(result.ReceiverId),
		IsRobot:    int32(result.IsRobot),
		UpdatedAt:  timeutil.DateTime(),
	}

	if item.TalkType == entity.ChatPrivateMode {
		item.UnreadNum = int32(c.unreadTalkCache.Get(ctx.RequestContext(), params.ReceiverId, uid))
		item.RemarkName = c.contactService.Dao().GetFriendRemark(ctx.RequestContext(), uid, params.ReceiverId, true)

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
	if msg, err := c.lastMessage.Get(ctx.RequestContext(), result.TalkType, uid, result.ReceiverId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return ctx.Success(&web.CreateTalkListResponse{
		Item: item,
	})
}

// Delete 删除列表
func (c *Talk) Delete(ctx *ichat.Context) error {

	params := &web.DeleteTalkListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.talkListService.Delete(ctx.Context, ctx.UserId(), params.Id); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(&web.DeleteTalkListResponse{})
}

// Top 置顶列表
func (c *Talk) Top(ctx *ichat.Context) error {

	params := &web.TopTalkListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.talkListService.Top(ctx.Context, &service.TalkSessionTopOpt{
		UserId: ctx.UserId(),
		Id:     params.Id,
		Type:   params.Type,
	}); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(&web.TopTalkListResponse{})
}

// Disturb 会话免打扰
func (c *Talk) Disturb(ctx *ichat.Context) error {

	params := &web.DisturbTalkListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.talkListService.Disturb(ctx.Context, &service.TalkSessionDisturbOpt{
		UserId:     ctx.UserId(),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		IsDisturb:  params.IsDisturb,
	}); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(&web.DisturbTalkListResponse{})
}

func (c *Talk) ClearUnreadMessage(ctx *ichat.Context) error {

	params := &web.ClearTalkUnreadNumRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if params.TalkType == 1 {
		c.unreadTalkCache.Reset(ctx.RequestContext(), params.ReceiverId, uid)
	}

	return ctx.Success(&web.ClearTalkUnreadNumResponse{})
}
