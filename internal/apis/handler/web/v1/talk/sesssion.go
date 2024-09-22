package talk

import (
	"fmt"
	"github.com/samber/lo"
	"go-chat/internal/repository/model"
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
		uid   = ctx.UserId()
		agent = strings.TrimSpace(ctx.Context.GetHeader("user-agent"))
	)

	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if agent != "" {
		agent = encrypt.Md5(agent)
	}

	// 判断对方是否是自己
	if in.TalkMode == entity.ChatPrivateMode && int(in.ToFromId) == ctx.UserId() {
		return ctx.ErrorBusiness("创建失败")
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, in.ToFromId, in.TalkMode, agent)
	if !c.RedisLock.Lock(ctx.Ctx(), key, 10) {
		return ctx.ErrorBusiness("创建失败")
	}

	if c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType: int(in.TalkMode),
		UserId:   uid,
		ToFromId: int(in.ToFromId),
	}) != nil {
		return ctx.ErrorBusiness("暂无权限！")
	}

	result, err := c.TalkSessionService.Create(ctx.Ctx(), &service.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   int(in.TalkMode),
		ReceiverId: int(in.ToFromId),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	item := &web.TalkSessionItem{
		Id:        int32(result.Id),
		TalkMode:  int32(result.TalkMode),
		ToFromId:  int32(result.ToFromId),
		IsTop:     int32(result.IsTop),
		IsDisturb: int32(result.IsDisturb),
		IsOnline:  model.No,
		IsRobot:   int32(result.IsRobot),
		Name:      "",
		Avatar:    "",
		Remark:    "",
		UnreadNum: 0,
		MsgText:   "",
		UpdatedAt: timeutil.DateTime(),
	}

	if item.TalkMode == entity.ChatPrivateMode {
		item.UnreadNum = int32(c.UnreadStorage.Get(ctx.Ctx(), uid, 1, int(in.ToFromId)))

		item.Remark = c.ContactRepo.GetFriendRemark(ctx.Ctx(), uid, int(in.ToFromId))
		if user, err := c.UsersRepo.FindById(ctx.Ctx(), result.ToFromId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if result.TalkMode == entity.ChatGroupMode {
		if group, err := c.GroupRepo.FindById(ctx.Ctx(), int(in.ToFromId)); err == nil {
			item.Name = group.Name
			item.Avatar = group.Avatar
		}
	}

	// 查询缓存消息
	if msg, err := c.MessageStorage.Get(ctx.Ctx(), result.TalkMode, uid, result.ToFromId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return ctx.Success(&web.TalkSessionCreateResponse{
		Id:        item.Id,
		TalkMode:  item.TalkMode,
		ToFromId:  item.ToFromId,
		IsTop:     item.IsTop,
		IsDisturb: item.IsDisturb,
		IsOnline:  item.IsOnline,
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
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Delete(ctx.Ctx(), ctx.UserId(), int(in.TalkMode), int(in.ToFromId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.TalkSessionDeleteResponse{})
}

// Top 置顶列表
func (c *Session) Top(ctx *core.Context) error {
	in := &web.TalkSessionTopRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Top(ctx.Ctx(), &service.TalkSessionTopOpt{
		UserId:   ctx.UserId(),
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		Action:   int(in.Action),
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.TalkSessionTopResponse{})
}

// Disturb 会话免打扰
func (c *Session) Disturb(ctx *core.Context) error {
	in := &web.TalkSessionDisturbRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkSessionService.Disturb(ctx.Ctx(), &service.TalkSessionDisturbOpt{
		UserId:   ctx.UserId(),
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		Action:   int(in.Action),
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.TalkSessionDisturbResponse{})
}

// List 会话列表
func (c *Session) List(ctx *core.Context) error {
	uid := ctx.UserId()

	data, err := c.TalkSessionService.List(ctx.Ctx(), uid)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkMode == 1 {
			friends = append(friends, item.ToFromId)
		}
	}

	// 获取好友备注
	remarks, _ := c.ContactRepo.Remarks(ctx.Ctx(), uid, friends)

	items := make([]*web.TalkSessionItem, 0)
	for _, item := range data {
		value := &web.TalkSessionItem{
			Id:        int32(item.Id),
			TalkMode:  int32(item.TalkMode),
			ToFromId:  int32(item.ToFromId),
			IsTop:     int32(item.IsTop),
			IsDisturb: int32(item.IsDisturb),
			IsRobot:   int32(item.IsRobot),
			IsOnline:  2,
			Avatar:    item.Avatar,
			MsgText:   "...",
			UpdatedAt: timeutil.FormatDatetime(item.UpdatedAt),
			UnreadNum: int32(c.UnreadStorage.Get(ctx.Ctx(), uid, item.TalkMode, item.ToFromId)),
		}

		if item.TalkMode == 1 {
			isOnline, _ := c.ClientConnectService.IsUidOnline(ctx.Ctx(), entity.ImChannelChat, int(value.ToFromId))

			value.Name = item.Nickname
			value.Avatar = item.Avatar
			value.Remark = remarks[item.ToFromId]
			value.IsOnline = lo.Ternary[int32](isOnline, 1, 2)
		} else {
			value.Name = item.GroupName
			value.Avatar = item.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := c.MessageStorage.Get(ctx.Ctx(), item.TalkMode, uid, item.ToFromId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	return ctx.Success(&web.TalkSessionListResponse{Items: items})
}

func (c *Session) ClearUnreadMessage(ctx *core.Context) error {
	in := &web.TalkSessionClearUnreadNumRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	c.UnreadStorage.Reset(ctx.Ctx(), ctx.UserId(), int(in.TalkMode), int(in.ToFromId))

	return ctx.Success(&web.TalkSessionClearUnreadNumResponse{})
}
