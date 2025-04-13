package talk

import (
	"fmt"
	"strings"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Session struct {
	RedisLock            *cache.RedisLock
	MessageStorage       *cache.MessageStorage
	ClientStorage        *cache.ClientStorage
	UnreadStorage        *cache.UnreadStorage
	ContactRemark        *cache.ContactRemark
	ContactRepo          *repo.Contact
	UsersRepo            *repo.Users
	GroupRepo            *repo.Group
	TalkService          service.ITalkService
	TalkSessionService   service.ITalkSessionService
	UserService          service.IUserService
	GroupService         service.IGroupService
	AuthService          service.IAuthService
	ContactService       service.IContactService
	ClientConnectService service.IClientConnectService
}

// Create 创建会话列表
func (c *Session) Create(ctx *core.Context) error {
	var (
		in    = &web.TalkSessionCreateRequest{}
		uid   = ctx.GetAuthId()
		agent = strings.TrimSpace(ctx.Context.GetHeader("user-agent"))
	)

	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if agent != "" {
		agent = encrypt.Md5(agent)
	}

	// 判断对方是否是自己
	if in.TalkMode == entity.ChatPrivateMode && int(in.ToFromId) == ctx.GetAuthId() {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, in.ToFromId, in.TalkMode, agent)
	if !c.RedisLock.Lock(ctx.GetContext(), key, 10) {
		return ctx.Error(entity.ErrTooFrequentOperation)
	}

	if c.AuthService.IsAuth(ctx.GetContext(), &service.AuthOption{
		TalkType: int(in.TalkMode),
		UserId:   uid,
		ToFromId: int(in.ToFromId),
	}) != nil {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	result, err := c.TalkSessionService.Create(ctx.GetContext(), &service.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   int(in.TalkMode),
		ReceiverId: int(in.ToFromId),
	})
	if err != nil {
		return ctx.Error(err)
	}

	item := &web.TalkSessionItem{
		Id:        int32(result.Id),
		TalkMode:  int32(result.TalkMode),
		ToFromId:  int32(result.ToFromId),
		IsTop:     int32(result.IsTop),
		IsDisturb: int32(result.IsDisturb),
		IsRobot:   int32(result.IsRobot),
		Name:      "",
		Avatar:    "",
		Remark:    "",
		UnreadNum: 0,
		MsgText:   "",
		UpdatedAt: timeutil.DateTime(),
	}

	if item.TalkMode == entity.ChatPrivateMode {
		item.UnreadNum = int32(c.UnreadStorage.Get(ctx.GetContext(), uid, 1, int(in.ToFromId)))

		item.Remark = c.ContactRepo.GetFriendRemark(ctx.GetContext(), uid, int(in.ToFromId))
		if user, err := c.UsersRepo.FindById(ctx.GetContext(), result.ToFromId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if result.TalkMode == entity.ChatGroupMode {
		if group, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.ToFromId)); err == nil {
			item.Name = group.Name
			item.Avatar = group.Avatar
		}
	}

	// 查询缓存消息
	if msg, err := c.MessageStorage.Get(ctx.GetContext(), result.TalkMode, uid, result.ToFromId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return ctx.Success(&web.TalkSessionCreateResponse{
		Id:        item.Id,
		TalkMode:  item.TalkMode,
		ToFromId:  item.ToFromId,
		IsTop:     item.IsTop,
		IsDisturb: item.IsDisturb,
		IsRobot:   item.IsRobot,
		Name:      item.Name,
		Avatar:    item.Avatar,
		Remark:    item.Remark,
		UnreadNum: item.UnreadNum,
		MsgText:   item.MsgText,
		UpdatedAt: item.UpdatedAt,
	})
}

// Delete 删除列表
func (c *Session) Delete(ctx *core.Context) error {
	in := &web.TalkSessionDeleteRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Delete(ctx.GetContext(), ctx.GetAuthId(), int(in.TalkMode), int(in.ToFromId)); err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.TalkSessionDeleteResponse{})
}

// Top 置顶列表
func (c *Session) Top(ctx *core.Context) error {
	in := &web.TalkSessionTopRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Top(ctx.GetContext(), &service.TalkSessionTopOpt{
		UserId:   ctx.GetAuthId(),
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		Action:   int(in.Action),
	}); err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.TalkSessionTopResponse{})
}

// Disturb 会话免打扰
func (c *Session) Disturb(ctx *core.Context) error {
	in := &web.TalkSessionDisturbRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Disturb(ctx.GetContext(), &service.TalkSessionDisturbOpt{
		UserId:   ctx.GetAuthId(),
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		Action:   int(in.Action),
	}); err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.TalkSessionDisturbResponse{})
}

// List 会话列表
func (c *Session) List(ctx *core.Context) error {
	uid := ctx.GetAuthId()

	data, err := c.TalkSessionService.List(ctx.GetContext(), uid)
	if err != nil {
		return ctx.Error(err)
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkMode == 1 {
			friends = append(friends, item.ToFromId)
		}
	}

	// 获取好友备注
	remarks, _ := c.ContactRepo.Remarks(ctx.GetContext(), uid, friends)

	items := make([]*web.TalkSessionItem, 0)
	for _, item := range data {
		value := &web.TalkSessionItem{
			Id:        int32(item.Id),
			TalkMode:  int32(item.TalkMode),
			ToFromId:  int32(item.ToFromId),
			IsTop:     int32(item.IsTop),
			IsDisturb: int32(item.IsDisturb),
			IsRobot:   int32(item.IsRobot),
			Avatar:    item.Avatar,
			MsgText:   "...",
			UpdatedAt: timeutil.FormatDatetime(item.UpdatedAt),
			UnreadNum: int32(c.UnreadStorage.Get(ctx.GetContext(), uid, item.TalkMode, item.ToFromId)),
		}

		if item.TalkMode == entity.ChatPrivateMode {
			value.Name = item.Nickname
			value.Avatar = item.Avatar
			value.Remark = remarks[item.ToFromId]
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := c.MessageStorage.Get(ctx.GetContext(), item.TalkMode, uid, item.ToFromId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	return ctx.Success(&web.TalkSessionListResponse{Items: items})
}

func (c *Session) ClearUnreadMessage(ctx *core.Context) error {
	in := &web.TalkSessionClearUnreadNumRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	c.UnreadStorage.Reset(ctx.GetContext(), ctx.GetAuthId(), int(in.TalkMode), int(in.ToFromId))

	return ctx.Success(&web.TalkSessionClearUnreadNumResponse{})
}
