package contact

import (
	"errors"

	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	message2 "go-chat/internal/service/message"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/service"
)

type Contact struct {
	ContactRepo     *repo.Contact
	UsersRepo       *repo.Users
	OrganizeRepo    *repo.Organize
	TalkSessionRepo *repo.TalkSession
	ContactService  service.IContactService
	UserService     service.IUserService
	TalkListService service.ITalkSessionService
	Message         message2.IService
	UserClient      *cache.UserClient
}

// List 联系人列表
func (c *Contact) List(ctx *core.Context) error {
	list, err := c.ContactService.List(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return ctx.Error(err)
	}

	items := make([]*web.ContactListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ContactListResponse_Item{
			UserId:   int32(item.Id),
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
func (c *Contact) Delete(ctx *core.Context) error {
	in := &web.ContactDeleteRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()
	if err := c.ContactService.Delete(ctx.GetContext(), uid, int(in.UserId)); err != nil {
		return ctx.Error(err)
	}

	_ = c.Message.CreatePrivateSysMessage(ctx.GetContext(), message2.CreatePrivateSysMessageOption{
		FromId:   int(in.UserId),
		ToFromId: uid,
		Content:  "你与对方已经解除了好友关系！",
	})

	if err := c.TalkListService.Delete(ctx.GetContext(), ctx.AuthId(), entity.ChatPrivateMode, int(in.UserId)); err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.ContactDeleteResponse{})
}

// Search 查找联系人
func (c *Contact) Search(ctx *core.Context) error {
	in := &web.ContactSearchRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.UsersRepo.FindByMobile(ctx.GetContext(), in.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Error(entity.ErrUserNotExist)
		}

		return ctx.Error(err)
	}

	return ctx.Success(&web.ContactSearchResponse{
		UserId:   int32(user.Id),
		Mobile:   lo.FromPtr[string](user.Mobile),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
	})
}

// Remark 编辑联系人备注
func (c *Contact) Remark(ctx *core.Context) error {
	in := &web.ContactEditRemarkRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ContactService.UpdateRemark(ctx.GetContext(), ctx.AuthId(), int(in.UserId), in.Remark); err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.ContactEditRemarkResponse{})
}

// Detail 联系人详情信息
func (c *Contact) Detail(ctx *core.Context) error {
	in := &web.ContactDetailRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	user, err := c.UsersRepo.FindByIdWithCache(ctx.GetContext(), int(in.UserId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Error(entity.ErrUserNotExist)
		}

		return ctx.Error(err)
	}

	resp := web.ContactDetailResponse{
		UserId:         int32(user.Id),
		Mobile:         lo.FromPtr(user.Mobile),
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Gender:         int32(user.Gender),
		Motto:          user.Motto,
		Email:          user.Email,
		Relation:       1, // 关系 1陌生人 2好友 3企业同事 4本人
		ContactRemark:  "",
		ContactGroupId: 0,
		OnlineStatus:   "N",
	}

	if user.Id == uid {
		resp.Relation = 4
		resp.OnlineStatus = "Y"
		return ctx.Success(&resp)
	}

	isQiYeMember, _ := c.OrganizeRepo.IsQiyeMember(ctx.GetContext(), uid, user.Id)
	if isQiYeMember {
		if c.UserClient.IsOnline(ctx.GetContext(), int64(in.UserId)) {
			resp.OnlineStatus = "Y"
		}

		resp.Relation = 3
		return ctx.Success(&resp)
	}

	contact, err := c.ContactRepo.FindByWhere(ctx.GetContext(), "user_id = ? and friend_id = ?", uid, user.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	resp.Relation = 1
	if err == nil && contact.Status == 1 && c.ContactRepo.IsFriend(ctx.GetContext(), uid, user.Id, true) {
		resp.Relation = 2
		resp.ContactGroupId = int32(contact.GroupId)
		resp.ContactRemark = contact.Remark

		if c.UserClient.IsOnline(ctx.GetContext(), int64(in.UserId)) {
			resp.OnlineStatus = "Y"
		}
	}

	return ctx.Success(&resp)
}

// MoveGroup 移动好友分组
func (c *Contact) MoveGroup(ctx *core.Context) error {
	in := &web.ContactChangeGroupRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ContactService.MoveGroup(ctx.GetContext(), ctx.AuthId(), int(in.UserId), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.ContactChangeGroupResponse{})
}

// OnlineStatus 获取联系人在线状态
func (c *Contact) OnlineStatus(ctx *core.Context) error {
	in := &web.ContactOnlineStatusRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	resp := &web.ContactOnlineStatusResponse{
		OnlineStatus: "N",
	}

	ok := c.ContactRepo.IsFriend(ctx.GetContext(), ctx.AuthId(), int(in.UserId), true)
	if ok && c.UserClient.IsOnline(ctx.GetContext(), int64(in.UserId)) {
		resp.OnlineStatus = "Y"
	}

	return ctx.Success(resp)
}
