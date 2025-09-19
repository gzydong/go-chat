package contact

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	message2 "go-chat/internal/service/message"
	"gorm.io/gorm"

	"go-chat/internal/entity"
	"go-chat/internal/service"
)

var _ web.IContactHandler = (*Contact)(nil)

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

func (c *Contact) List(ctx context.Context, _ *web.ContactListRequest) (*web.ContactListResponse, error) {
	list, err := c.ContactService.List(ctx, middleware.FormContextAuthId[entity.WebClaims](ctx))
	if err != nil {
		return nil, err
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

	return &web.ContactListResponse{Items: items}, nil
}

func (c *Contact) Delete(ctx context.Context, in *web.ContactDeleteRequest) (*web.ContactDeleteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if err := c.ContactService.Delete(ctx, uid, int(in.UserId)); err != nil {
		return nil, err
	}

	_ = c.Message.CreatePrivateSysMessage(ctx, message2.CreatePrivateSysMessageOption{
		FromId:   int(in.UserId),
		ToFromId: uid,
		Content:  "你与对方已经解除了好友关系！",
	})

	if err := c.TalkListService.Delete(ctx, uid, entity.ChatPrivateMode, int(in.UserId)); err != nil {
		return nil, err
	}

	return &web.ContactDeleteResponse{}, nil
}

func (c *Contact) EditRemark(ctx context.Context, in *web.ContactEditRemarkRequest) (*web.ContactEditRemarkResponse, error) {
	if err := c.ContactService.UpdateRemark(ctx, middleware.FormContextAuthId[entity.WebClaims](ctx), int(in.UserId), in.Remark); err != nil {
		return nil, err
	}

	return &web.ContactEditRemarkResponse{}, nil
}

func (c *Contact) Detail(ctx context.Context, in *web.ContactDetailRequest) (*web.ContactDetailResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	user, err := c.UsersRepo.FindByIdWithCache(ctx, int(in.UserId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotExist
		}

		return nil, err
	}

	resp := &web.ContactDetailResponse{
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
		return resp, nil
	}

	isQiYeMember, _ := c.OrganizeRepo.IsQiyeMember(ctx, uid, user.Id)
	if isQiYeMember {
		if c.UserClient.IsOnline(ctx, int64(in.UserId)) {
			resp.OnlineStatus = "Y"
		}

		resp.Relation = 3
		return resp, nil
	}

	contact, err := c.ContactRepo.FindByWhere(ctx, "user_id = ? and friend_id = ?", uid, user.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	resp.Relation = 1
	if err == nil && contact.Status == 1 && c.ContactRepo.IsFriend(ctx, uid, user.Id, true) {
		resp.Relation = 2
		resp.ContactGroupId = int32(contact.GroupId)
		resp.ContactRemark = contact.Remark

		if c.UserClient.IsOnline(ctx, int64(in.UserId)) {
			resp.OnlineStatus = "Y"
		}
	}

	return resp, nil
}

func (c *Contact) Search(ctx context.Context, in *web.ContactSearchRequest) (*web.ContactSearchResponse, error) {
	user, err := c.UsersRepo.FindByMobile(ctx, in.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUserNotExist
		}

		return nil, err
	}

	return &web.ContactSearchResponse{
		UserId:   int32(user.Id),
		Mobile:   lo.FromPtr[string](user.Mobile),
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Gender:   int32(user.Gender),
		Motto:    user.Motto,
	}, nil
}

func (c *Contact) ChangeGroup(ctx context.Context, in *web.ContactChangeGroupRequest) (*web.ContactChangeGroupResponse, error) {
	err := c.ContactService.MoveGroup(ctx, middleware.FormContextAuthId[entity.WebClaims](ctx), int(in.UserId), int(in.GroupId))
	if err != nil {
		return nil, err
	}

	return &web.ContactChangeGroupResponse{}, nil
}

func (c *Contact) OnlineStatus(ctx context.Context, in *web.ContactOnlineStatusRequest) (*web.ContactOnlineStatusResponse, error) {
	resp := &web.ContactOnlineStatusResponse{
		OnlineStatus: "N",
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	ok := c.ContactRepo.IsFriend(ctx, uid, int(in.UserId), true)
	if ok && c.UserClient.IsOnline(ctx, int64(in.UserId)) {
		resp.OnlineStatus = "Y"
	}

	return resp, nil
}
