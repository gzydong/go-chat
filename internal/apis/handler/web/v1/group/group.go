package group

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
)

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

// Create 创建群聊分组
func (c *Group) Create(ctx *core.Context) error {
	in := &web.GroupCreateRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) < 2 {
		return ctx.InvalidParams("创建群聊失败，至少需要两个用户！")
	}

	if len(uids)+1 > model.GroupMemberMaxNum {
		return ctx.InvalidParams(fmt.Sprintf("群成员数量已达到%d上限！", model.GroupMemberMaxNum))
	}

	gid, err := c.GroupService.Create(ctx.GetContext(), &service.GroupCreateOpt{
		UserId:    ctx.AuthId(),
		Name:      in.Name,
		MemberIds: uids,
	})

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.GroupCreateResponse{GroupId: int32(gid)})
}

// Dismiss 解散群组
func (c *Group) Dismiss(ctx *core.Context) error {
	in := &web.GroupDismissRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsMaster(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	if err := c.GroupService.Dismiss(ctx.GetContext(), int(in.GroupId), uid); err != nil {
		return ctx.Error(err)
	}

	_ = c.Message.CreateGroupSysMessage(ctx.GetContext(), message.CreateGroupSysMessageOption{
		GroupId: int(in.GroupId),
		Content: "该群已被群主解散！",
	})

	return ctx.Success(nil)
}

// Invite 邀请好友加入群聊
func (c *Group) Invite(ctx *core.Context) error {
	in := &web.GroupInviteRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) == 0 {
		return ctx.InvalidParams("邀请好友列表不能为空！")
	}

	if len(uids) > model.GroupMemberMaxNum {
		return ctx.InvalidParams(fmt.Sprintf("当前群成员数量已达到%d上限！", model.GroupMemberMaxNum))
	}

	key := fmt.Sprintf("group_join:%d", in.GroupId)
	if !c.RedisLock.Lock(ctx.GetContext(), key, 20) {
		return ctx.Error(entity.ErrTooFrequentOperation)
	}

	defer c.RedisLock.UnLock(ctx.GetContext(), key)

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsMember(ctx.GetContext(), int(in.GroupId), uid, true) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	group, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	if group != nil && group.IsDismiss == model.Yes {
		return ctx.Error(entity.ErrGroupDismissed)
	}

	count, err := c.GroupMemberRepo.FindCount(ctx.GetContext(), "group_id = ? and is_quit = ?", in.GroupId, model.No)
	if err != nil {
		return ctx.Error(err)
	}

	if int(count)+len(uids) >= model.GroupMemberMaxNum {
		return ctx.Error(entity.ErrGroupMemberLimit)
	}

	if err := c.GroupService.Invite(ctx.GetContext(), &service.GroupInviteOpt{
		UserId:    uid,
		GroupId:   int(in.GroupId),
		MemberIds: uids,
	}); err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.GroupInviteResponse{})
}

// Secede 退出群聊
func (c *Group) Secede(ctx *core.Context) error {
	in := &web.GroupSecedeRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()
	if err := c.GroupService.Secede(ctx.GetContext(), int(in.GroupId), uid); err != nil {
		return ctx.Error(err)
	}

	_ = c.TalkSessionService.Delete(ctx.GetContext(), uid, entity.ChatGroupMode, int(in.GroupId))

	return ctx.Success(nil)
}

// Update 群设置接口（预留）
func (c *Group) Update(ctx *core.Context) error {
	in := &web.GroupSettingRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	group, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	if group != nil && group.IsDismiss == model.Yes {
		return ctx.Error(entity.ErrGroupDismissed)
	}

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsLeader(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	if err := c.GroupService.Update(ctx.GetContext(), &service.GroupUpdateOpt{
		GroupId: int(in.GroupId),
		Name:    in.GroupName,
		Avatar:  in.Avatar,
		Profile: in.Profile,
	}); err != nil {
		return ctx.Error(err)
	}

	_ = c.Message.CreateGroupSysMessage(ctx.GetContext(), message.CreateGroupSysMessageOption{
		GroupId: int(in.GroupId),
		Content: "群主或管理员修改了群信息！",
	})

	return ctx.Success(&web.GroupSettingResponse{})
}

// RemoveMember 移除指定成员(群组&管理员权限)
func (c *Group) RemoveMember(ctx *core.Context) error {
	in := &web.GroupRemoveMemberRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uids := make([]int, 0)
	for _, id := range sliceutil.Unique(in.UserIds) {
		uids = append(uids, int(id))
	}

	if len(uids) == 0 {
		return ctx.InvalidParams("移除成员列表不能为空！")
	}

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsLeader(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	err := c.GroupService.RemoveMember(ctx.GetContext(), &service.GroupRemoveMembersOpt{
		UserId:    uid,
		GroupId:   int(in.GroupId),
		MemberIds: uids,
	})

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.GroupRemoveMemberResponse{})
}

// Detail 获取群组信息
func (c *Group) Detail(ctx *core.Context) error {
	in := &web.GroupDetailRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	groupInfo, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	if groupInfo.Id == 0 {
		return ctx.Error(entity.ErrGroupNotExist)
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
		VisitCard: c.GroupMemberRepo.GetMemberRemark(ctx.GetContext(), int(in.GroupId), uid),
		Notice: &web.GroupDetailResponse_Notice{
			Content:        "",
			CreatedAt:      "",
			UpdatedAt:      "",
			ModifyUserName: "",
		},
	}

	notice, err := c.GroupNoticeRepo.GetLatestNotice(ctx.GetContext(), int(in.GroupId))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if notice != nil {
		resp.Notice = &web.GroupDetailResponse_Notice{
			Content:        notice.Content,
			CreatedAt:      timeutil.FormatDatetime(notice.CreatedAt),
			UpdatedAt:      timeutil.FormatDatetime(notice.UpdatedAt),
			ModifyUserName: "",
		}
	}

	if c.TalkSessionRepo.IsDisturb(uid, groupInfo.Id, 2) {
		resp.IsDisturb = 1
	}

	return ctx.Success(resp)
}

// UpdateMemberRemark 修改群备注接口
func (c *Group) UpdateMemberRemark(ctx *core.Context) error {
	in := &web.GroupRemarkUpdateRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := c.GroupMemberRepo.UpdateByWhere(ctx.GetContext(), map[string]any{
		"user_card": in.Remark,
	}, "group_id = ? and user_id = ?", in.GroupId, ctx.AuthId())
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

func (c *Group) GetInviteFriends(ctx *core.Context) error {
	in := &web.GetInviteFriendsRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	items, err := c.ContactService.List(ctx.GetContext(), ctx.AuthId())
	if err != nil {
		return ctx.Error(err)
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

		return ctx.Success(&web.GetInviteFriendsResponse{
			Items: data,
		})
	}

	mids := c.GroupMemberRepo.GetMemberIds(ctx.GetContext(), int(in.GroupId))
	if len(mids) == 0 {
		return ctx.Success(&web.GetInviteFriendsResponse{
			Items: data,
		})
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

	return ctx.Success(&web.GetInviteFriendsResponse{
		Items: data,
	})
}

func (c *Group) List(ctx *core.Context) error {
	items, err := c.GroupService.List(ctx.AuthId())
	if err != nil {
		return ctx.Error(err)
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

	return ctx.Success(resp)
}

// Members 获取群成员列表
func (c *Group) Members(ctx *core.Context) error {
	in := &web.GroupMemberListRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	group, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	if group != nil && group.IsDismiss == model.Yes {
		return ctx.Success(&web.GroupMemberListResponse{})
	}

	if !c.GroupMemberRepo.IsMember(ctx.GetContext(), int(in.GroupId), ctx.AuthId(), false) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	list := c.GroupMemberRepo.GetMembers(ctx.GetContext(), int(in.GroupId))

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

	return ctx.Success(&web.GroupMemberListResponse{Items: items})
}

// OvertList 公开群列表
func (c *Group) OvertList(ctx *core.Context) error {
	in := &web.GroupOvertListRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	list, err := c.GroupRepo.SearchOvertList(ctx.GetContext(), &repo.SearchOvertListOpt{
		Name:   in.Name,
		UserId: uid,
		Page:   int(in.Page),
		Size:   20,
	})
	if err != nil {
		return ctx.Error(err)
	}

	resp := &web.GroupOvertListResponse{}
	resp.Items = make([]*web.GroupOvertListResponse_Item, 0)

	if len(list) == 0 {
		return ctx.Success(resp)
	}

	ids := make([]int, 0)
	for _, val := range list {
		ids = append(ids, val.Id)
	}

	count, err := c.GroupMemberRepo.CountGroupMemberNum(ids)
	if err != nil {
		return ctx.Error(err)
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

	return ctx.Success(resp)
}

// Transfer 群主转让
func (c *Group) Transfer(ctx *core.Context) error {
	in := &web.GroupHandoverRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsMaster(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	if uid == int(in.UserId) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	err := c.GroupMemberService.Handover(ctx.GetContext(), int(in.GroupId), uid, int(in.UserId))
	if err != nil {
		return ctx.Error(err)
	}

	members := make([]model.TalkRecordExtraGroupMember, 0)
	c.Repo.Db().Table("users").Select("id as user_id", "nickname").Where("id in ?", []int{uid, int(in.UserId)}).Scan(&members)

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

	_ = c.Message.CreateGroupMessage(ctx.GetContext(), message.CreateGroupMessageOption{
		MsgType:  entity.ChatMsgSysGroupTransfer,
		FromId:   uid,
		ToFromId: int(in.GroupId),
		Extra:    jsonutil.Encode(extra),
	})

	return ctx.Success(nil)
}

// AssignAdmin 分配管理员
func (c *Group) AssignAdmin(ctx *core.Context) error {
	in := &web.GroupAssignAdminRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsMaster(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	leader := lo.Ternary(in.Action == 1, model.GroupMemberLeaderAdmin, model.GroupMemberLeaderOrdinary)

	err := c.GroupMemberService.SetLeaderStatus(ctx.GetContext(), int(in.GroupId), int(in.UserId), leader)
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

// MemberMute 禁止发言
func (c *Group) MemberMute(ctx *core.Context) error {
	in := &web.GroupNoSpeakRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()
	if !c.GroupMemberRepo.IsLeader(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	status := lo.Ternary(in.Action == 1, model.Yes, model.No)

	err := c.GroupMemberService.SetMuteStatus(ctx.GetContext(), int(in.GroupId), int(in.UserId), status)
	if err != nil {
		return ctx.Error(err)
	}

	members := make([]model.TalkRecordExtraGroupMember, 0)
	c.Repo.Db().Model(&model.Users{}).Select("id as user_id", "nickname").Where("id = ?", in.UserId).Scan(&members)

	user, err := c.UsersRepo.FindByIdWithCache(ctx.GetContext(), uid)
	if err != nil {
		return ctx.Error(err)
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

	_ = c.Message.CreateGroupMessage(ctx.GetContext(), data)

	return ctx.Success(nil)
}

// Mute 全员禁言
func (c *Group) Mute(ctx *core.Context) error {
	in := &web.GroupMuteRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	group, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	if group.IsDismiss == model.Yes {
		return ctx.Error(entity.ErrGroupDismissed)
	}

	if !c.GroupMemberRepo.IsLeader(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	data := map[string]any{
		"is_mute":    in.Action,
		"updated_at": time.Now(),
	}

	affected, err := c.GroupRepo.UpdateByWhere(ctx.GetContext(), data, "id = ?", in.GroupId)
	if err != nil {
		return ctx.Error(err)
	}

	if affected == 0 {
		return ctx.Success(web.GroupMuteResponse{})
	}

	user, err := c.UsersRepo.FindById(ctx.GetContext(), uid)
	if err != nil {
		return err
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

	_ = c.Message.CreateGroupMessage(ctx.GetContext(), message.CreateGroupMessageOption{
		MsgType:  msgType,
		FromId:   uid,
		ToFromId: int(in.GroupId),
		Extra:    jsonutil.Encode(extra),
	})

	return ctx.Success(web.GroupMuteResponse{})
}

// Overt 公开群
func (c *Group) Overt(ctx *core.Context) error {
	in := &web.GroupOvertRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	group, err := c.GroupRepo.FindById(ctx.GetContext(), int(in.GroupId))
	if err != nil {
		return ctx.Error(err)
	}

	if group.IsDismiss == model.Yes {
		return ctx.Error(entity.ErrGroupDismissed)
	}

	if !c.GroupMemberRepo.IsMaster(ctx.GetContext(), int(in.GroupId), uid) {
		return ctx.Error(entity.ErrPermissionDenied)
	}

	_, err = c.GroupRepo.UpdateByWhere(ctx.GetContext(), map[string]any{
		"is_overt":   in.Action,
		"updated_at": time.Now(),
	}, "id = ?", in.GroupId)

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(web.GroupOvertResponse{})
}
