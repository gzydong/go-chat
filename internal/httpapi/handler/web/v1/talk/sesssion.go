package talk

import (
	"fmt"
	"strconv"
	"strings"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Session struct {
	RedisLock          *cache.RedisLock
	MessageStorage     *cache.MessageStorage
	ClientStorage      *cache.ClientStorage
	UnreadStorage      *cache.UnreadStorage
	ContactRemark      *cache.ContactRemark
	ContactRepo        *repo.Contact
	UsersRepo          *repo.Users
	GroupRepo          *repo.Group
	TalkService        service.ITalkService
	TalkSessionService service.ITalkSessionService
	UserService        service.IUserService
	GroupService       service.IGroupService
	AuthService        service.IAuthService
	ContactService     service.IContactService
}

// Create 创建会话列表
func (c *Session) Create(ctx *ichat.Context) error {

	var (
		params = &web.TalkSessionCreateRequest{}
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
	if params.TalkType == entity.ChatPrivateMode && int(params.ReceiverId) == ctx.UserId() {
		return ctx.ErrorBusiness("创建失败")
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, params.ReceiverId, params.TalkType, agent)
	if !c.RedisLock.Lock(ctx.Ctx(), key, 10) {
		return ctx.ErrorBusiness("创建失败")
	}

	if c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   int(params.TalkType),
		UserId:     uid,
		ReceiverId: int(params.ReceiverId),
	}) != nil {
		return ctx.ErrorBusiness("暂无权限！")
	}

	result, err := c.TalkSessionService.Create(ctx.Ctx(), &service.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   int(params.TalkType),
		ReceiverId: int(params.ReceiverId),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	item := &web.TalkSessionItem{
		Id:         int32(result.Id),
		TalkType:   int32(result.TalkType),
		ReceiverId: int32(result.ReceiverId),
		IsRobot:    int32(result.IsRobot),
		UpdatedAt:  timeutil.DateTime(),
	}

	if item.TalkType == entity.ChatPrivateMode {
		item.UnreadNum = int32(c.UnreadStorage.Get(ctx.Ctx(), 1, int(params.ReceiverId), uid))
		item.Remark = c.ContactRepo.GetFriendRemark(ctx.Ctx(), uid, int(params.ReceiverId))

		if user, err := c.UsersRepo.FindById(ctx.Ctx(), result.ReceiverId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if result.TalkType == entity.ChatGroupMode {
		if group, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.ReceiverId)); err == nil {
			item.Name = group.Name
		}
	}

	// 查询缓存消息
	if msg, err := c.MessageStorage.Get(ctx.Ctx(), result.TalkType, uid, result.ReceiverId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return ctx.Success(&web.TalkSessionCreateResponse{
		Id:         item.Id,
		TalkType:   item.TalkType,
		ReceiverId: item.ReceiverId,
		IsTop:      item.IsTop,
		IsDisturb:  item.IsDisturb,
		IsOnline:   item.IsOnline,
		IsRobot:    item.IsRobot,
		Name:       item.Name,
		Avatar:     item.Avatar,
		RemarkName: item.Remark,
		UnreadNum:  item.UnreadNum,
		MsgText:    item.MsgText,
		UpdatedAt:  item.UpdatedAt,
	})
}

// Delete 删除列表
func (c *Session) Delete(ctx *ichat.Context) error {

	params := &web.TalkSessionDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Delete(ctx.Ctx(), ctx.UserId(), int(params.ListId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.TalkSessionDeleteResponse{})
}

// Top 置顶列表
func (c *Session) Top(ctx *ichat.Context) error {

	params := &web.TalkSessionTopRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Top(ctx.Ctx(), &service.TalkSessionTopOpt{
		UserId: ctx.UserId(),
		Id:     int(params.ListId),
		Type:   int(params.Type),
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.TalkSessionTopResponse{})
}

// Disturb 会话免打扰
func (c *Session) Disturb(ctx *ichat.Context) error {

	params := &web.TalkSessionDisturbRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Disturb(ctx.Ctx(), &service.TalkSessionDisturbOpt{
		UserId:     ctx.UserId(),
		TalkType:   int(params.TalkType),
		ReceiverId: int(params.ReceiverId),
		IsDisturb:  int(params.IsDisturb),
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.TalkSessionDisturbResponse{})
}

// List 会话列表
func (c *Session) List(ctx *ichat.Context) error {

	uid := ctx.UserId()

	// 获取未读消息数
	unReads := c.UnreadStorage.All(ctx.Ctx(), uid)
	if len(unReads) > 0 {
		c.TalkSessionService.BatchAddList(ctx.Ctx(), uid, unReads)
	}

	data, err := c.TalkSessionService.List(ctx.Ctx(), uid)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkType == 1 {
			friends = append(friends, item.ReceiverId)
		}
	}

	// 获取好友备注
	remarks, _ := c.ContactRepo.Remarks(ctx.Ctx(), uid, friends)

	items := make([]*web.TalkSessionItem, 0)
	for _, item := range data {
		value := &web.TalkSessionItem{
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

		if num, ok := unReads[fmt.Sprintf("%d_%d", item.TalkType, item.ReceiverId)]; ok {
			value.UnreadNum = int32(num)
		}

		if item.TalkType == 1 {
			value.Name = item.Nickname
			value.Avatar = item.UserAvatar
			value.Remark = remarks[item.ReceiverId]
			value.IsOnline = int32(strutil.BoolToInt(c.ClientStorage.IsOnline(ctx.Ctx(), entity.ImChannelChat, strconv.Itoa(int(value.ReceiverId)))))
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := c.MessageStorage.Get(ctx.Ctx(), item.TalkType, uid, item.ReceiverId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	return ctx.Success(&web.TalkSessionListResponse{Items: items})
}

func (c *Session) ClearUnreadMessage(ctx *ichat.Context) error {

	params := &web.TalkSessionClearUnreadNumRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	c.UnreadStorage.Reset(ctx.Ctx(), int(params.TalkType), int(params.ReceiverId), ctx.UserId())

	return ctx.Success(&web.TalkSessionClearUnreadNumResponse{})
}
