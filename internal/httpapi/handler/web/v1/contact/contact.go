package contact

import (
	"errors"

	"go-chat/api/pb/message/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/service"
)

type Contact struct {
	ClientStorage   *cache.ClientStorage
	ContactRepo     *repo.Contact
	UsersRepo       *repo.Users
	OrganizeRepo    *repo.Organize
	TalkSessionRepo *repo.TalkSession
	ContactService  service.IContactService
	UserService     service.IUserService
	TalkListService service.ITalkSessionService
	MessageService  service.IMessageService
}

// List 联系人列表
func (c *Contact) List(ctx *ichat.Context) error {

	list, err := c.ContactService.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
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
			GroupId:  int32(item.GroupId),
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
	if err := c.ContactService.Delete(ctx.Ctx(), uid, int(params.FriendId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	// 解除好友关系后需添加一条聊天记录
	_ = c.MessageService.SendSystemText(ctx.Ctx(), uid, &message.TextMessageRequest{
		Content: "你与对方已经解除了好友关系！",
		Receiver: &message.MessageReceiver{
			TalkType:   entity.ChatPrivateMode,
			ReceiverId: params.FriendId,
		},
	})

	// 删除聊天会话
	sid := c.TalkSessionRepo.FindBySessionId(uid, int(params.FriendId), entity.ChatPrivateMode)
	if err := c.TalkListService.Delete(ctx.Ctx(), ctx.UserId(), sid); err != nil {
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

	user, err := c.UsersRepo.FindByMobile(params.Mobile)
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

// Remark 编辑联系人备注
func (c *Contact) Remark(ctx *ichat.Context) error {

	params := &web.ContactEditRemarkRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ContactService.UpdateRemark(ctx.Ctx(), ctx.UserId(), int(params.FriendId), params.Remark); err != nil {
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

	user, err := c.UsersRepo.FindById(ctx.Ctx(), int(params.UserId))
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

		contact, err := c.ContactRepo.FindByWhere(ctx.Ctx(), "user_id = ? and friend_id = ?", uid, user.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil && contact.Status == 1 {
			if c.ContactRepo.IsFriend(ctx.Ctx(), uid, user.Id, false) {
				data.FriendStatus = 2
				data.GroupId = int32(contact.GroupId)
				data.Remark = contact.Remark
			}
		} else {
			isOk, _ := c.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid, user.Id)
			if isOk {
				data.FriendStatus = 2
			}
		}
	}

	return ctx.Success(&data)
}

// MoveGroup 移动好友分组
func (c *Contact) MoveGroup(ctx *ichat.Context) error {

	params := &web.ContactChangeGroupRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ContactService.MoveGroup(ctx.Ctx(), ctx.UserId(), int(params.UserId), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactChangeGroupResponse{})
}
