package group

import (
	"context"
	"errors"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/logic"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/jsonutil"
	"github.com/gzydong/go-chat/internal/pkg/sliceutil"
	"github.com/gzydong/go-chat/internal/pkg/timeutil"
	"github.com/gzydong/go-chat/internal/repository/cache"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var _ web.IGroupApplyHandler = (*Apply)(nil)

type Apply struct {
	Redis              *redis.Client
	GroupApplyStorage  *cache.GroupApplyStorage
	GroupRepo          *repo.Group
	GroupApplyRepo     *repo.GroupApply
	GroupMemberRepo    *repo.GroupMember
	GroupApplyService  service.IGroupApplyService
	GroupMemberService service.IGroupMemberService
	GroupService       service.IGroupService
	PushMessage        *logic.PushMessage
}

func (a Apply) Create(ctx context.Context, in *web.GroupApplyCreateRequest) (*web.GroupApplyCreateResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	apply, err := a.GroupApplyRepo.FindByWhere(ctx, "group_id = ? and user_id = ? and status = ?", in.GroupId, uid, model.GroupApplyStatusWait)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	applyId := 0
	if apply == nil {
		data := &model.GroupApply{
			GroupId: int(in.GroupId),
			UserId:  uid,
			Status:  model.GroupApplyStatusWait,
			Remark:  in.Remark,
		}

		err = a.GroupApplyRepo.Create(ctx, data)
		if err == nil {
			applyId = data.Id
		}
	} else {
		applyId = apply.Id
		data := map[string]any{
			"remark":     in.Remark,
			"updated_at": timeutil.DateTime(),
		}

		_, err = a.GroupApplyRepo.UpdateByWhere(ctx, data, "id = ?", apply.Id)
	}

	if err != nil {
		return nil, err
	}

	find, err := a.GroupMemberRepo.FindByWhere(ctx, "group_id = ? and leader = ?", in.GroupId, model.GroupMemberLeaderOwner)
	if err == nil && find != nil {
		a.GroupApplyStorage.Incr(ctx, find.UserId)
	}

	_ = a.PushMessage.Push(ctx, entity.ImChannelChat, &entity.SubscribeMessage{
		Event: entity.SubEventGroupApply,
		Payload: jsonutil.Encode(entity.SubEventGroupApplyPayload{
			GroupId: int(in.GroupId),
			UserId:  uid,
			ApplyId: applyId,
		}),
	})

	return &web.GroupApplyCreateResponse{}, err
}

func (a Apply) Delete(ctx context.Context, req *web.GroupApplyDeleteRequest) (*web.GroupApplyDeleteResponse, error) {
	return nil, nil
}

func (a Apply) Agree(ctx context.Context, in *web.GroupApplyAgreeRequest) (*web.GroupApplyAgreeResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	apply, err := a.GroupApplyRepo.FindById(ctx, int(in.ApplyId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrDataNotFound
	}

	if !a.GroupMemberRepo.IsLeader(ctx, apply.GroupId, uid) {
		return nil, entity.ErrPermissionDenied
	}

	if apply.Status != model.GroupApplyStatusWait {
		return nil, nil
	}

	if !a.GroupMemberRepo.IsMember(ctx, apply.GroupId, apply.UserId, false) {
		err = a.GroupService.Invite(ctx, &service.GroupInviteOpt{
			UserId:    uid,
			GroupId:   apply.GroupId,
			MemberIds: []int{apply.UserId},
		})

		if err != nil {
			return nil, err
		}
	}

	data := map[string]any{
		"status":     model.GroupApplyStatusPass,
		"updated_at": timeutil.DateTime(),
	}

	_, err = a.GroupApplyRepo.UpdateByWhere(ctx, data, "id = ?", in.ApplyId)
	if err != nil {
		return nil, err
	}

	return &web.GroupApplyAgreeResponse{}, nil
}

func (a Apply) Decline(ctx context.Context, in *web.GroupApplyDeclineRequest) (*web.GroupApplyDeclineResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	apply, err := a.GroupApplyRepo.FindById(ctx, int(in.ApplyId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrDataNotFound
	}

	if !a.GroupMemberRepo.IsLeader(ctx, apply.GroupId, uid) {
		return nil, entity.ErrPermissionDenied
	}

	if apply.Status != model.GroupApplyStatusWait {
		return &web.GroupApplyDeclineResponse{}, nil
	}

	data := map[string]any{
		"status":     model.GroupApplyStatusRefuse,
		"reason":     in.Remark,
		"updated_at": timeutil.DateTime(),
	}

	_, err = a.GroupApplyRepo.UpdateByWhere(ctx, data, "id = ?", in.ApplyId)
	if err != nil {
		return nil, err
	}

	return &web.GroupApplyDeclineResponse{}, nil
}

func (a Apply) List(ctx context.Context, in *web.GroupApplyListRequest) (*web.GroupApplyListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if !a.GroupMemberRepo.IsLeader(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	list, err := a.GroupApplyRepo.List(ctx, []int{int(in.GroupId)})
	if err != nil {
		return nil, err
	}

	items := make([]*web.GroupApplyListResponse_Item, 0)
	for _, item := range list {
		items = append(items, &web.GroupApplyListResponse_Item{
			Id:        int32(item.Id),
			UserId:    int32(item.UserId),
			GroupId:   int32(item.GroupId),
			Remark:    item.Remark,
			Avatar:    item.Avatar,
			Nickname:  item.Nickname,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	return &web.GroupApplyListResponse{Items: items}, nil
}

func (a Apply) All(ctx context.Context, req *web.GroupApplyAllRequest) (*web.GroupApplyAllResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	all, err := a.GroupMemberRepo.FindAll(ctx, func(db *gorm.DB) {
		db.Select("group_id")
		db.Where("user_id = ?", uid)
		db.Where("leader in ?", []int{
			model.GroupMemberLeaderOwner,
			model.GroupMemberLeaderAdmin,
		})
		db.Where("is_quit = ?", model.No)
	})

	if err != nil {
		return nil, err
	}

	groupIds := make([]int, 0, len(all))
	for _, m := range all {
		groupIds = append(groupIds, m.GroupId)
	}

	resp := &web.GroupApplyAllResponse{Items: make([]*web.GroupApplyAllResponse_Item, 0)}

	if len(groupIds) == 0 {
		a.GroupApplyStorage.Del(ctx, uid)
		return resp, nil
	}

	list, err := a.GroupApplyRepo.List(ctx, groupIds)
	if err != nil {
		return nil, err
	}

	groups, err := a.GroupRepo.FindAll(ctx, func(db *gorm.DB) {
		db.Select("id,name")
		db.Where("id in ?", groupIds)
	})
	if err != nil {
		return nil, err
	}

	groupMap := sliceutil.ToMap(groups, func(t *model.Group) int {
		return t.Id
	})

	for _, item := range list {
		resp.Items = append(resp.Items, &web.GroupApplyAllResponse_Item{
			Id:        int32(item.Id),
			UserId:    int32(item.UserId),
			GroupName: groupMap[item.GroupId].Name,
			GroupId:   int32(item.GroupId),
			Remark:    item.Remark,
			Avatar:    item.Avatar,
			Nickname:  item.Nickname,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	a.GroupApplyStorage.Del(ctx, uid)

	return resp, nil
}

func (a Apply) UnreadNum(ctx context.Context, req *web.GroupApplyUnreadNumRequest) (*web.GroupApplyUnreadNumResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	return &web.GroupApplyUnreadNumResponse{
		Num: int32(a.GroupApplyStorage.Get(ctx, uid)),
	}, nil
}
