package contact

import (
	"context"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
)

var _ web.IContactApplyHandler = (*Apply)(nil)

type Apply struct {
	ContactRepo         *repo.Contact
	ContactApplyService service.IContactApplyService
	UserService         service.IUserService
	ContactService      service.IContactService
	MessageService      message.IService
}

func (a Apply) Create(ctx context.Context, in *web.ContactApplyCreateRequest) (*web.ContactApplyCreateResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if a.ContactRepo.IsFriend(ctx, uid, int(in.UserId), false) {
		return nil, nil
	}

	if err := a.ContactApplyService.Create(ctx, &service.ContactApplyCreateOpt{
		UserId:   uid,
		Remarks:  in.Remark,
		FriendId: int(in.UserId),
	}); err != nil {
		return nil, err
	}

	return &web.ContactApplyCreateResponse{}, nil
}

func (a Apply) Accept(ctx context.Context, in *web.ContactApplyAcceptRequest) (*web.ContactApplyAcceptResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	applyInfo, err := a.ContactApplyService.Accept(ctx, &service.ContactApplyAcceptOpt{
		Remarks: in.Remark,
		ApplyId: int(in.ApplyId),
		UserId:  uid,
	})

	if err != nil {
		return nil, err
	}

	_ = a.MessageService.CreatePrivateSysMessage(ctx, message.CreatePrivateSysMessageOption{
		FromId:   uid,
		ToFromId: applyInfo.UserId,
		Content:  "你们已成为好友，可以开始聊天咯！",
	})

	_ = a.MessageService.CreatePrivateSysMessage(ctx, message.CreatePrivateSysMessageOption{
		FromId:   applyInfo.UserId,
		ToFromId: uid,
		Content:  "你们已成为好友，可以开始聊天咯！",
	})

	return &web.ContactApplyAcceptResponse{}, nil
}

func (a Apply) Decline(ctx context.Context, in *web.ContactApplyDeclineRequest) (*web.ContactApplyDeclineResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := a.ContactApplyService.Decline(ctx, &service.ContactApplyDeclineOpt{
		UserId:  uid,
		Remarks: in.Remark,
		ApplyId: int(in.ApplyId),
	}); err != nil {
		return nil, err
	}

	return &web.ContactApplyDeclineResponse{}, nil
}

func (a Apply) List(ctx context.Context, req *web.ContactApplyListRequest) (*web.ContactApplyListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	list, err := a.ContactApplyService.List(ctx, uid)
	if err != nil {
		return nil, err
	}

	items := make([]*web.ContactApplyListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ContactApplyListResponse_Item{
			Id:        int32(item.Id),
			UserId:    int32(item.UserId),
			FriendId:  int32(item.FriendId),
			Remark:    item.Remark,
			Nickname:  item.Nickname,
			Avatar:    item.Avatar,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	a.ContactApplyService.ClearApplyUnreadNum(ctx, uid)

	return &web.ContactApplyListResponse{Items: items}, nil
}

func (a Apply) UnreadNum(ctx context.Context, req *web.ContactApplyUnreadNumRequest) (*web.ContactApplyUnreadNumResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	return &web.ContactApplyUnreadNumResponse{Num: int32(a.ContactApplyService.GetApplyUnreadNum(ctx, uid))}, nil
}
