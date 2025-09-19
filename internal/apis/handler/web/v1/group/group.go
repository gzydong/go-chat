package group

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/errorx"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
	"gorm.io/gorm"
)

var _ web.IGroupHandler = (*Group)(nil)

type Group struct {
	RedisLock          *cache.RedisLock
	Repo               *repo.Source
	UsersRepo          *repo.Users
	GroupRepo          *repo.Group
	GroupMemberRepo    *repo.GroupMember
	GroupNoticeRepo    *repo.GroupNotice
	TalkSessionRepo    *repo.TalkSession
	GroupService       service.IGroupService
	GroupMemberService service.IGroupMemberService
	TalkSessionService service.ITalkSessionService
	UserService        service.IUserService
	ContactService     service.IContactService
	Message            message.IService
}

func (g Group) List(ctx context.Context, in *web.GroupListRequest) (*web.GroupListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	items, err := g.GroupService.List(uid)
	if err != nil {
		return nil, err
	}

	resp := &web.GroupListResponse{
		Items: make([]*web.GroupListResponse_Item, 0, len(items)),
	}

	for _, item := range items {
		resp.Items = append(resp.Items, &web.GroupListResponse_Item{
			GroupId:   int32(item.Id),
			GroupName: item.GroupName,
			Avatar:    item.Avatar,
			Profile:   item.Profile,
			Leader:    int32(item.Leader),
			CreatorId: int32(item.CreatorId),
		})
	}

	return resp, nil
}

func (g Group) Create(ctx context.Context, in *web.GroupCreateRequest) (*web.GroupCreateResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) < 2 {
		return nil, errorx.New(400, "创建群聊失败，至少需要两个用户")
	}

	if len(uids)+1 > model.GroupMemberMaxNum {
		return nil, errorx.New(400, fmt.Sprintf("群成员数量已达到%d上限！", model.GroupMemberMaxNum))
	}

	gid, err := g.GroupService.Create(ctx, &service.GroupCreateOpt{
		UserId:    uid,
		Name:      in.Name,
		MemberIds: uids,
	})

	if err != nil {
		return nil, err
	}

	return &web.GroupCreateResponse{GroupId: int32(gid)}, nil
}

func (g Group) Detail(ctx context.Context, in *web.GroupDetailRequest) (*web.GroupDetailResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	groupInfo, err := g.GroupRepo.FindById(ctx, int(in.GroupId))
	if err != nil {
		return nil, err
	}

	if groupInfo.Id == 0 {
		return nil, entity.ErrGroupNotExist
	}

	resp := &web.GroupDetailResponse{
		GroupId:   int32(groupInfo.Id),
		GroupName: groupInfo.Name,
		Profile:   groupInfo.Profile,
		Avatar:    groupInfo.Avatar,
		CreatedAt: timeutil.FormatDatetime(groupInfo.CreatedAt),
		IsManager: uid == groupInfo.CreatorId,
		IsDisturb: 0,
		IsMute:    int32(groupInfo.IsMute),
		IsOvert:   int32(groupInfo.IsOvert),
		VisitCard: g.GroupMemberRepo.GetMemberRemark(ctx, int(in.GroupId), uid),
		Notice: &web.GroupDetailResponse_Notice{
			Content:        "",
			CreatedAt:      "",
			UpdatedAt:      "",
			ModifyUserName: "",
		},
	}

	notice, err := g.GroupNoticeRepo.GetLatestNotice(ctx, int(in.GroupId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if notice != nil {
		resp.Notice = &web.GroupDetailResponse_Notice{
			Content:        notice.Content,
			CreatedAt:      timeutil.FormatDatetime(notice.CreatedAt),
			UpdatedAt:      timeutil.FormatDatetime(notice.UpdatedAt),
			ModifyUserName: "",
		}
	}

	if g.TalkSessionRepo.IsDisturb(uid, groupInfo.Id, 2) {
		resp.IsDisturb = 1
	}

	return resp, nil
}

func (g Group) MemberList(ctx context.Context, in *web.GroupMemberListRequest) (*web.GroupMemberListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	group, err := g.GroupRepo.FindById(ctx, int(in.GroupId))
	if err != nil {
		return nil, err
	}

	if group != nil && group.IsDismiss == model.Yes {
		return &web.GroupMemberListResponse{}, nil
	}

	if !g.GroupMemberRepo.IsMember(ctx, int(in.GroupId), uid, false) {
		return nil, entity.ErrPermissionDenied
	}

	list := g.GroupMemberRepo.GetMembers(ctx, int(in.GroupId))

	items := make([]*web.GroupMemberListResponse_Item, 0)
	for _, item := range list {
		items = append(items, &web.GroupMemberListResponse_Item{
			UserId:   int32(item.UserId),
			Nickname: item.Nickname,
			Avatar:   item.Avatar,
			Gender:   int32(item.Gender),
			Leader:   int32(item.Leader),
			IsMute:   int32(item.IsMute),
			Remark:   item.UserCard,
			Motto:    item.Motto,
		})
	}

	slices.SortFunc(items, func(a, b *web.GroupMemberListResponse_Item) int {
		return int(a.Leader - b.Leader)
	})

	return &web.GroupMemberListResponse{Items: items}, nil
}

func (g Group) Dismiss(ctx context.Context, in *web.GroupDismissRequest) (*web.GroupDismissResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if !g.GroupMemberRepo.IsMaster(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	if err := g.GroupService.Dismiss(ctx, int(in.GroupId), uid); err != nil {
		return nil, err
	}

	_ = g.Message.CreateGroupSysMessage(ctx, message.CreateGroupSysMessageOption{
		GroupId: int(in.GroupId),
		Content: "该群已被群主解散！",
	})

	return &web.GroupDismissResponse{}, nil
}

func (g Group) Invite(ctx context.Context, in *web.GroupInviteRequest) (*web.GroupInviteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) == 0 {
		return nil, errorx.New(400, "邀请好友列表不能为空")
	}

	if len(uids) > model.GroupMemberMaxNum {
		return nil, errorx.New(400, fmt.Sprintf("当前群成员数量已达到%d上限！", model.GroupMemberMaxNum))
	}

	key := fmt.Sprintf("group_join:%d", in.GroupId)
	if !g.RedisLock.Lock(ctx, key, 20) {
		return nil, entity.ErrTooFrequentOperation
	}

	defer g.RedisLock.UnLock(ctx, key)

	if !g.GroupMemberRepo.IsMember(ctx, int(in.GroupId), uid, true) {
		return nil, entity.ErrPermissionDenied
	}

	group, err := g.GroupRepo.FindById(ctx, int(in.GroupId))
	if err != nil {
		return nil, err
	}

	if group != nil && group.IsDismiss == model.Yes {
		return nil, entity.ErrGroupDismissed
	}

	count, err := g.GroupMemberRepo.FindCount(ctx, "group_id = ? and is_quit = ?", in.GroupId, model.No)
	if err != nil {
		return nil, err
	}

	if int(count)+len(uids) >= model.GroupMemberMaxNum {
		return nil, entity.ErrGroupMemberLimit
	}

	if err := g.GroupService.Invite(ctx, &service.GroupInviteOpt{
		UserId:    uid,
		GroupId:   int(in.GroupId),
		MemberIds: uids,
	}); err != nil {
		return nil, err
	}

	return &web.GroupInviteResponse{}, nil
}

func (g Group) GetInviteFriends(ctx context.Context, in *web.GetInviteFriendsRequest) (*web.GetInviteFriendsResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	items, err := g.ContactService.List(ctx, uid)
	if err != nil {
		return nil, err
	}

	data := make([]*web.GetInviteFriendsResponse_Item, 0)
	if in.GroupId <= 0 {
		for _, item := range items {
			data = append(data, &web.GetInviteFriendsResponse_Item{
				UserId:   int32(item.Id),
				Nickname: item.Nickname,
				Avatar:   item.Avatar,
				Gender:   int32(item.Gender),
				Remark:   item.Remark,
			})
		}

		return &web.GetInviteFriendsResponse{
			Items: data,
		}, nil
	}

	mids := g.GroupMemberRepo.GetMemberIds(ctx, int(in.GroupId))
	if len(mids) == 0 {
		return &web.GetInviteFriendsResponse{
			Items: data,
		}, nil
	}

	for i := 0; i < len(items); i++ {
		if !slices.Contains(mids, items[i].Id) {
			data = append(data, &web.GetInviteFriendsResponse_Item{
				UserId:   int32(items[i].Id),
				Nickname: items[i].Nickname,
				Avatar:   items[i].Avatar,
				Gender:   int32(items[i].Gender),
				Remark:   items[i].Remark,
			})
		}
	}

	return &web.GetInviteFriendsResponse{
		Items: data,
	}, nil
}

func (g Group) Secede(ctx context.Context, in *web.GroupSecedeRequest) (*web.GroupSecedeResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := g.GroupService.Secede(ctx, int(in.GroupId), uid); err != nil {
		return nil, err
	}

	_ = g.TalkSessionService.Delete(ctx, uid, entity.ChatGroupMode, int(in.GroupId))

	return &web.GroupSecedeResponse{}, nil
}

func (g Group) Setting(ctx context.Context, req *web.GroupSettingRequest) (*web.GroupSettingResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g Group) RemarkUpdate(ctx context.Context, in *web.GroupRemarkUpdateRequest) (*web.GroupRemarkUpdateResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	_, err := g.GroupMemberRepo.UpdateByWhere(ctx, map[string]any{
		"user_card": in.Remark,
	}, "group_id = ? and user_id = ?", in.GroupId, uid)
	if err != nil {
		return nil, err
	}

	return &web.GroupRemarkUpdateResponse{}, nil
}

func (g Group) RemoveMember(ctx context.Context, in *web.GroupRemoveMemberRequest) (*web.GroupRemoveMemberResponse, error) {
	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) == 0 {
		return nil, errorx.New(400, "移除成员列表不能为空！")
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if !g.GroupMemberRepo.IsLeader(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	err := g.GroupService.RemoveMember(ctx, &service.GroupRemoveMembersOpt{
		UserId:    uid,
		GroupId:   int(in.GroupId),
		MemberIds: uids,
	})

	if err != nil {
		return nil, err
	}

	return &web.GroupRemoveMemberResponse{}, nil
}

func (g Group) OvertList(ctx context.Context, in *web.GroupOvertListRequest) (*web.GroupOvertListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	list, err := g.GroupRepo.SearchOvertList(ctx, &repo.SearchOvertListOpt{
		Name:   in.Name,
		UserId: uid,
		Page:   int(in.Page),
		Size:   20,
	})
	if err != nil {
		return nil, err
	}

	resp := &web.GroupOvertListResponse{}
	resp.Items = make([]*web.GroupOvertListResponse_Item, 0)

	if len(list) == 0 {
		return resp, nil
	}

	ids := make([]int, 0)
	for _, val := range list {
		ids = append(ids, val.Id)
	}

	count, err := g.GroupMemberRepo.CountGroupMemberNum(ids)
	if err != nil {
		return nil, err
	}

	countMap := make(map[int]int)
	for _, member := range count {
		countMap[member.GroupId] = member.Count
	}

	for i, value := range list {
		if i >= 19 {
			break
		}

		resp.Items = append(resp.Items, &web.GroupOvertListResponse_Item{
			GroupId:   int32(value.Id),
			Type:      int32(value.Type),
			Name:      value.Name,
			Avatar:    value.Avatar,
			Profile:   value.Profile,
			Count:     int32(countMap[value.Id]),
			MaxNum:    int32(value.MaxNum),
			CreatedAt: timeutil.FormatDatetime(value.CreatedAt),
		})
	}

	resp.Next = len(list) > 19

	return resp, nil
}

func (g Group) Handover(ctx context.Context, in *web.GroupHandoverRequest) (*web.GroupHandoverResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if !g.GroupMemberRepo.IsMaster(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	if uid == int(in.UserId) {
		return nil, entity.ErrPermissionDenied
	}

	err := g.GroupMemberService.Handover(ctx, int(in.GroupId), uid, int(in.UserId))
	if err != nil {
		return nil, err
	}

	members := make([]model.TalkRecordExtraGroupMember, 0)
	g.Repo.Db().Table("users").Select("id as user_id", "nickname").Where("id in ?", []int{uid, int(in.UserId)}).Scan(&members)

	extra := model.TalkRecordExtraTransferGroup{}
	for _, member := range members {
		if member.UserId == uid {
			extra.OldOwnerId = member.UserId
			extra.OldOwnerName = member.Nickname
		} else {
			extra.NewOwnerId = member.UserId
			extra.NewOwnerName = member.Nickname
		}
	}

	_ = g.Message.CreateGroupMessage(ctx, message.CreateGroupMessageOption{
		MsgType:  entity.ChatMsgSysGroupTransfer,
		FromId:   uid,
		ToFromId: int(in.GroupId),
		Extra:    jsonutil.Encode(extra),
	})

	return &web.GroupHandoverResponse{}, nil
}

func (g Group) AssignAdmin(ctx context.Context, in *web.GroupAssignAdminRequest) (*web.GroupAssignAdminResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if !g.GroupMemberRepo.IsMaster(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	leader := lo.Ternary(in.Action == 1, model.GroupMemberLeaderAdmin, model.GroupMemberLeaderOrdinary)

	err := g.GroupMemberService.SetLeaderStatus(ctx, int(in.GroupId), int(in.UserId), leader)
	if err != nil {
		return nil, err
	}

	return &web.GroupAssignAdminResponse{}, nil
}

func (g Group) NoSpeak(ctx context.Context, in *web.GroupNoSpeakRequest) (*web.GroupNoSpeakResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if !g.GroupMemberRepo.IsLeader(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	status := lo.Ternary(in.Action == 1, model.Yes, model.No)

	err := g.GroupMemberService.SetMuteStatus(ctx, int(in.GroupId), int(in.UserId), status)
	if err != nil {
		return nil, err
	}

	members := make([]model.TalkRecordExtraGroupMember, 0)
	g.Repo.Db().Model(&model.Users{}).Select("id as user_id", "nickname").Where("id = ?", in.UserId).Scan(&members)

	user, err := g.UsersRepo.FindByIdWithCache(ctx, uid)
	if err != nil {
		return nil, err
	}

	data := message.CreateGroupMessageOption{
		FromId:   uid,
		ToFromId: int(in.GroupId),
	}

	if status == model.Yes {
		data.MsgType = entity.ChatMsgSysGroupMemberMuted
		data.Extra = jsonutil.Encode(model.TalkRecordExtraGroupMemberCancelMuted{
			OwnerId:   uid,
			OwnerName: user.Nickname,
			Members:   members,
		})
	} else {
		data.MsgType = entity.ChatMsgSysGroupMemberCancelMuted
		data.Extra = jsonutil.Encode(model.TalkRecordExtraGroupMemberCancelMuted{
			OwnerId:   uid,
			OwnerName: user.Nickname,
			Members:   members,
		})
	}

	_ = g.Message.CreateGroupMessage(ctx, data)

	return &web.GroupNoSpeakResponse{}, nil
}

func (g Group) Mute(ctx context.Context, in *web.GroupMuteRequest) (*web.GroupMuteResponse, error) {

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	group, err := g.GroupRepo.FindById(ctx, int(in.GroupId))
	if err != nil {
		return nil, err
	}

	if group.IsDismiss == model.Yes {
		return nil, entity.ErrGroupDismissed
	}

	if !g.GroupMemberRepo.IsLeader(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	data := map[string]any{
		"is_mute":    in.Action,
		"updated_at": time.Now(),
	}

	affected, err := g.GroupRepo.UpdateByWhere(ctx, data, "id = ?", in.GroupId)
	if err != nil {
		return nil, err
	}

	if affected == 0 {
		return &web.GroupMuteResponse{}, nil
	}

	user, err := g.UsersRepo.FindById(ctx, uid)
	if err != nil {
		return nil, err
	}

	var extra any
	var msgType int
	if in.Action == model.Yes {
		msgType = entity.ChatMsgSysGroupMuted
		extra = model.TalkRecordExtraGroupMuted{
			OwnerId:   user.Id,
			OwnerName: user.Nickname,
		}
	} else {
		msgType = entity.ChatMsgSysGroupCancelMuted
		extra = model.TalkRecordExtraGroupCancelMuted{
			OwnerId:   user.Id,
			OwnerName: user.Nickname,
		}
	}

	_ = g.Message.CreateGroupMessage(ctx, message.CreateGroupMessageOption{
		MsgType:  msgType,
		FromId:   uid,
		ToFromId: int(in.GroupId),
		Extra:    jsonutil.Encode(extra),
	})

	return &web.GroupMuteResponse{}, nil
}

func (g Group) Overt(ctx context.Context, in *web.GroupOvertRequest) (*web.GroupOvertResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	group, err := g.GroupRepo.FindById(ctx, int(in.GroupId))
	if err != nil {
		return nil, err
	}

	if group.IsDismiss == model.Yes {
		return nil, entity.ErrGroupDismissed
	}

	if !g.GroupMemberRepo.IsMaster(ctx, int(in.GroupId), uid) {
		return nil, entity.ErrPermissionDenied
	}

	_, err = g.GroupRepo.UpdateByWhere(ctx, map[string]any{
		"is_overt":   in.Action,
		"updated_at": time.Now(),
	}, "id = ?", in.GroupId)

	if err != nil {
		return nil, err
	}

	return &web.GroupOvertResponse{}, nil
}
