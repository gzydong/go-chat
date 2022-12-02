package contact

import (
	"errors"
	"strconv"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/cache"
	"go-chat/internal/service/organize"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

type Contact struct {
	service            *service.ContactService
	wsClient           *cache.ClientStorage
	userService        *service.UserService
	talkListService    *service.TalkSessionService
	talkMessageService *service.TalkMessageService
	organizeService    *organize.OrganizeService
}

func NewContact(service *service.ContactService, wsClient *cache.ClientStorage, userService *service.UserService, talkListService *service.TalkSessionService, talkMessageService *service.TalkMessageService, organizeService *organize.OrganizeService) *Contact {
	return &Contact{service: service, wsClient: wsClient, userService: userService, talkListService: talkListService, talkMessageService: talkMessageService, organizeService: organizeService}
}

// List 联系人列表
func (c *Contact) List(ctx *ichat.Context) error {

	list, err := c.service.List(ctx.Context, ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	for _, item := range list {
		item.IsOnline = strutil.BoolToInt(c.wsClient.IsOnline(ctx.Context, entity.ImChannelChat, strconv.Itoa(item.Id)))
	}

	items := make([]*web.ContactListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ContactListResponse_Item{
			Id:       int32(item.Id),
			Nickname: item.Nickname,
			Gender:   int32(item.Gender),
			Motto:    item.Motto,
			Avatar:   item.Avatar,
			Remark:   item.Remark,
			IsOnline: int32(strutil.BoolToInt(c.wsClient.IsOnline(ctx.Context, entity.ImChannelChat, strconv.Itoa(item.Id)))),
		})
	}

	return ctx.Success(&web.ContactListResponse{Items: items})
}

// Delete 删除联系人
func (c *Contact) Delete(ctx *ichat.Context) error {

	params := &web.ContactDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.service.Delete(ctx.Context, uid, int(params.FriendId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	// 解除好友关系后需添加一条聊天记录
	_ = c.talkMessageService.SendSysMessage(ctx.Context, &service.SysTextMessageOpt{
		UserId:     uid,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: int(params.FriendId),
		Text:       "你与对方已经解除了好友关系！！！",
	})

	// 删除聊天会话
	sid := c.talkListService.Dao().FindBySessionId(uid, int(params.FriendId), entity.ChatPrivateMode)
	if err := c.talkListService.Delete(ctx.Context, ctx.UserId(), sid); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactDeleteResponse{})
}

// Search 查找联系人
func (c *Contact) Search(ctx *ichat.Context) error {

	params := &web.ContactSearchRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.userService.Dao().FindByMobile(params.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.ErrorBusiness("用户不存在！")
		}

		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactSearchResponse{
		Id:       int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
	})
}

// EditRemark 编辑联系人备注
func (c *Contact) EditRemark(ctx *ichat.Context) error {

	params := &web.ContactEditRemarkRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.EditRemark(ctx.Context, ctx.UserId(), int(params.FriendId), params.Remark); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactEditRemarkResponse{})
}

// Detail 联系人详情信息
func (c *Contact) Detail(ctx *ichat.Context) error {

	params := &web.ContactDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	user, err := c.userService.Dao().FindById(int(params.UserId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.ErrorBusiness("用户不存在！")
		}

		return ctx.ErrorBusiness(err.Error())
	}

	data := web.ContactDetailResponse{
		Id:           int32(user.Id),
		Mobile:       user.Mobile,
		Nickname:     user.Nickname,
		Avatar:       user.Avatar,
		Gender:       int32(user.Gender),
		Motto:        user.Motto,
		FriendApply:  0,
		FriendStatus: 0, // 朋友关系[0:本人;1:陌生人;2:朋友;]
	}

	if uid != user.Id {
		data.FriendStatus = 1
		if c.service.Dao().IsFriend(ctx.Ctx(), uid, user.Id, false) {
			data.FriendStatus = 2
			data.NicknameRemark = c.service.Dao().GetFriendRemark(ctx.Ctx(), uid, user.Id)
		} else {
			isOk, _ := c.organizeService.Dao().IsQiyeMember(uid, user.Id)
			if isOk {
				data.FriendStatus = 2
			}
		}
	}

	return ctx.Success(&data)
}
