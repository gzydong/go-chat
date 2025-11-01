package talk

import (
	"context"
	"fmt"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/timeutil"
	"github.com/gzydong/go-chat/internal/repository/cache"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
)

var _ web.ITalkHandler = (*Session)(nil)

type Session struct {
	RedisLock          *cache.RedisLock
	MessageStorage     *cache.MessageStorage
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
}

func (s *Session) SessionCreate(ctx context.Context, in *web.TalkSessionCreateRequest) (*web.TalkSessionCreateResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	agent := "" // TODO aa
	// 判断对方是否是自己
	if in.TalkMode == entity.ChatPrivateMode && int(in.ToFromId) == uid {
		return nil, entity.ErrPermissionDenied
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d-%s", uid, in.ToFromId, in.TalkMode, agent)
	if !s.RedisLock.Lock(ctx, key, 10) {
		return nil, entity.ErrTooFrequentOperation
	}

	if s.AuthService.IsAuth(ctx, &service.AuthOption{
		TalkType: int(in.TalkMode),
		UserId:   uid,
		ToFromId: int(in.ToFromId),
	}) != nil {
		return nil, entity.ErrPermissionDenied
	}

	result, err := s.TalkSessionService.Create(ctx, &service.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   int(in.TalkMode),
		ReceiverId: int(in.ToFromId),
	})
	if err != nil {
		return nil, err
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
		item.UnreadNum = int32(s.UnreadStorage.Get(ctx, uid, 1, int(in.ToFromId)))

		item.Remark = s.ContactRepo.GetFriendRemark(ctx, uid, int(in.ToFromId))
		if user, err := s.UsersRepo.FindById(ctx, result.ToFromId); err == nil {
			item.Name = user.Nickname
			item.Avatar = user.Avatar
		}
	} else if result.TalkMode == entity.ChatGroupMode {
		if group, err := s.GroupRepo.FindById(ctx, int(in.ToFromId)); err == nil {
			item.Name = group.Name
			item.Avatar = group.Avatar
		}
	}

	// 查询缓存消息
	if msg, err := s.MessageStorage.Get(ctx, result.TalkMode, uid, result.ToFromId); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	return &web.TalkSessionCreateResponse{
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
	}, nil
}

func (s *Session) SessionDelete(ctx context.Context, in *web.TalkSessionDeleteRequest) (*web.TalkSessionDeleteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := s.TalkSessionService.Delete(ctx, uid, int(in.TalkMode), int(in.ToFromId)); err != nil {
		return nil, err
	}

	return &web.TalkSessionDeleteResponse{}, nil
}

func (s *Session) SessionTop(ctx context.Context, in *web.TalkSessionTopRequest) (*web.TalkSessionTopResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if err := s.TalkSessionService.Top(ctx, &service.TalkSessionTopOpt{
		UserId:   uid,
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		Action:   int(in.Action),
	}); err != nil {
		return nil, err
	}

	return &web.TalkSessionTopResponse{}, nil
}

func (s *Session) SessionDisturb(ctx context.Context, in *web.TalkSessionDisturbRequest) (*web.TalkSessionDisturbResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if err := s.TalkSessionService.Disturb(ctx, &service.TalkSessionDisturbOpt{
		UserId:   uid,
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		Action:   int(in.Action),
	}); err != nil {
		return nil, err
	}

	return &web.TalkSessionDisturbResponse{}, nil
}

func (s *Session) SessionList(ctx context.Context, req *web.TalkSessionListRequest) (*web.TalkSessionListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	data, err := s.TalkSessionService.List(ctx, uid)
	if err != nil {
		return nil, err
	}

	friends := make([]int, 0)
	for _, item := range data {
		if item.TalkMode == 1 {
			friends = append(friends, item.ToFromId)
		}
	}

	// 获取好友备注
	remarks, _ := s.ContactRepo.Remarks(ctx, uid, friends)

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
			UnreadNum: int32(s.UnreadStorage.Get(ctx, uid, item.TalkMode, item.ToFromId)),
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
		if msg, err := s.MessageStorage.Get(ctx, item.TalkMode, uid, item.ToFromId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}

		items = append(items, value)
	}

	return &web.TalkSessionListResponse{Items: items}, nil
}

func (s *Session) SessionClearUnreadNum(ctx context.Context, in *web.TalkSessionClearUnreadNumRequest) (*web.TalkSessionClearUnreadNumResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	s.UnreadStorage.Reset(ctx, uid, int(in.TalkMode), int(in.ToFromId))
	return &web.TalkSessionClearUnreadNumResponse{}, nil
}
