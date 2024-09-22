package contact

import (
	"errors"
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
	ClientStorage   *cache.ClientStorage
	ContactRepo     *repo.Contact
	UsersRepo       *repo.Users
	OrganizeRepo    *repo.Organize
	TalkSessionRepo *repo.TalkSession
	ContactService  service.IContactService
	UserService     service.IUserService
	TalkListService service.ITalkSessionService
	Message         message2.IService
}

// List 联系人列表
func (c *Contact) List(ctx *core.Context) error {
	list, err := c.ContactService.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
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
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.ContactService.Delete(ctx.Ctx(), uid, int(in.UserId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	_ = c.Message.CreatePrivateSysMessage(ctx.Ctx(), message2.CreatePrivateSysMessageOption{
		FromId:   int(in.UserId),
		ToFromId: uid,
		Content:  "你与对方已经解除了好友关系！",
	})

	if err := c.TalkListService.Delete(ctx.Ctx(), ctx.UserId(), entity.ChatPrivateMode, int(in.UserId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactDeleteResponse{})
}

// Search 查找联系人
func (c *Contact) Search(ctx *core.Context) error {
	in := &web.ContactSearchRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	user, err := c.UsersRepo.FindByMobile(ctx.Ctx(), in.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.ErrorBusiness("用户不存在！")
		}

		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactSearchResponse{
		UserId:   int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
	})
}

// Remark 编辑联系人备注
func (c *Contact) Remark(ctx *core.Context) error {
	in := &web.ContactEditRemarkRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ContactService.UpdateRemark(ctx.Ctx(), ctx.UserId(), int(in.UserId), in.Remark); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactEditRemarkResponse{})
}

// Detail 联系人详情信息
func (c *Contact) Detail(ctx *core.Context) error {
	in := &web.ContactDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	user, err := c.UsersRepo.FindByIdWithCache(ctx.Ctx(), int(in.UserId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.ErrorBusiness("用户不存在！")
		}

		return ctx.ErrorBusiness(err.Error())
	}

	data := web.ContactDetailResponse{
		UserId:   int32(user.Id),
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
		Email:    user.Email,
		FriendInfo: &web.ContactDetailResponse_FriendInfo{
			IsFriend: "N",
			GroupId:  0,
			Remark:   "",
		},
	}

	if uid != user.Id {
		contact, err := c.ContactRepo.FindByWhere(ctx.Ctx(), "user_id = ? and friend_id = ?", uid, user.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err == nil && contact.Status == 1 {
			if c.ContactRepo.IsFriend(ctx.Ctx(), uid, user.Id, false) {
				data.FriendInfo.IsFriend = "Y"
				data.FriendInfo.GroupId = int32(contact.GroupId)
				data.FriendInfo.Remark = contact.Remark
			}
		} else {
			isOk, _ := c.OrganizeRepo.IsQiyeMember(ctx.Ctx(), uid, user.Id)
			if isOk {
				data.FriendInfo.IsFriend = "Y"
			}
		}
	}

	return ctx.Success(&data)
}

// MoveGroup 移动好友分组
func (c *Contact) MoveGroup(ctx *core.Context) error {
	in := &web.ContactChangeGroupRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ContactService.MoveGroup(ctx.Ctx(), ctx.UserId(), int(in.UserId), int(in.GroupId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ContactChangeGroupResponse{})
}
