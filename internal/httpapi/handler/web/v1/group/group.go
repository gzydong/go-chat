package group

import (
	"fmt"
	"slices"

	"go-chat/api/pb/message/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Group struct {
	RedisLock          *cache.RedisLock
	Repo               *repo.Source
	UsersRepo          *repo.Users
	GroupRepo          *repo.Group
	GroupMemberRepo    *repo.GroupMember
	TalkSessionRepo    *repo.TalkSession
	GroupService       service.IGroupService
	GroupMemberService service.IGroupMemberService
	TalkSessionService service.ITalkSessionService
	UserService        service.IUserService
	ContactService     service.IContactService
	GroupNoticeService service.IGroupNoticeService
	MessageService     service.IMessageService
}

// Create 创建群聊分组
func (c *Group) Create(ctx *ichat.Context) error {

	params := &web.GroupCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	gid, err := c.GroupService.Create(ctx.Ctx(), &service.GroupCreateOpt{
		UserId:    ctx.UserId(),
		Name:      params.Name,
		Avatar:    params.Avatar,
		MemberIds: sliceutil.ParseIds(params.GetIds()),
	})
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！" + err.Error())
	}

	return ctx.Success(&web.GroupCreateResponse{GroupId: int32(gid)})
}

// Dismiss 解散群组
func (c *Group) Dismiss(ctx *ichat.Context) error {

	params := &web.GroupDismissRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.GroupMemberRepo.IsMaster(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限解散群组！")
	}

	if err := c.GroupService.Dismiss(ctx.Ctx(), int(params.GroupId), ctx.UserId()); err != nil {
		return ctx.ErrorBusiness("群组解散失败！")
	}

	_ = c.MessageService.SendSystemText(ctx.Ctx(), uid, &message.TextMessageRequest{
		Content: "群组已被群主解散！",
		Receiver: &message.MessageReceiver{
			TalkType:   entity.ChatGroupMode,
			ReceiverId: params.GroupId,
		},
	})

	return ctx.Success(nil)
}

// Invite 邀请好友加入群聊
func (c *Group) Invite(ctx *ichat.Context) error {

	params := &web.GroupInviteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	key := fmt.Sprintf("group-join:%d", params.GroupId)
	if !c.RedisLock.Lock(ctx.Ctx(), key, 20) {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	defer c.RedisLock.UnLock(ctx.Ctx(), key)

	group, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	if group != nil && group.IsDismiss == 1 {
		return ctx.ErrorBusiness("该群已解散！")
	}

	uid := ctx.UserId()
	uids := sliceutil.Unique(sliceutil.ParseIds(params.Ids))

	if len(uids) == 0 {
		return ctx.ErrorBusiness("邀请好友列表不能为空！")
	}

	if !c.GroupMemberRepo.IsMember(ctx.Ctx(), int(params.GroupId), uid, true) {
		return ctx.ErrorBusiness("非群组成员，无权邀请好友！")
	}

	if err := c.GroupService.Invite(ctx.Ctx(), &service.GroupInviteOpt{
		UserId:    uid,
		GroupId:   int(params.GroupId),
		MemberIds: uids,
	}); err != nil {
		return ctx.ErrorBusiness("邀请好友加入群聊失败！" + err.Error())
	}

	return ctx.Success(&web.GroupInviteResponse{})
}

// SignOut 退出群聊
func (c *Group) SignOut(ctx *ichat.Context) error {

	params := &web.GroupSecedeRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.GroupService.Secede(ctx.Ctx(), int(params.GroupId), uid); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	sid := c.TalkSessionRepo.FindBySessionId(uid, int(params.GroupId), entity.ChatGroupMode)

	_ = c.TalkSessionService.Delete(ctx.Ctx(), ctx.UserId(), sid)

	return ctx.Success(nil)
}

// Setting 群设置接口（预留）
func (c *Group) Setting(ctx *ichat.Context) error {

	params := &web.GroupSettingRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	group, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	if group != nil && group.IsDismiss == 1 {
		return ctx.ErrorBusiness("该群已解散！")
	}

	uid := ctx.UserId()
	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("无权限操作")
	}

	if err := c.GroupService.Update(ctx.Ctx(), &service.GroupUpdateOpt{
		GroupId: int(params.GroupId),
		Name:    params.GroupName,
		Avatar:  params.Avatar,
		Profile: params.Profile,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	_ = c.MessageService.SendSystemText(ctx.Ctx(), uid, &message.TextMessageRequest{
		Content: "群主或管理员修改了群信息！",
		Receiver: &message.MessageReceiver{
			TalkType:   entity.ChatGroupMode,
			ReceiverId: params.GroupId,
		},
	})

	return ctx.Success(&web.GroupSettingResponse{})
}

// RemoveMembers 移除指定成员(群组&管理员权限)
func (c *Group) RemoveMembers(ctx *ichat.Context) error {

	params := &web.GroupRemoveMemberRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("无权限操作")
	}

	err := c.GroupService.RemoveMember(ctx.Ctx(), &service.GroupRemoveMembersOpt{
		UserId:    uid,
		GroupId:   int(params.GroupId),
		MemberIds: sliceutil.ParseIds(params.MembersIds),
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.GroupRemoveMemberResponse{})
}

// Detail 获取群组信息
func (c *Group) Detail(ctx *ichat.Context) error {

	params := &web.GroupDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	groupInfo, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if groupInfo.Id == 0 {
		return ctx.ErrorBusiness("数据不存在")
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
		VisitCard: c.GroupMemberRepo.GetMemberRemark(ctx.Ctx(), int(params.GroupId), uid),
	}

	if c.TalkSessionRepo.IsDisturb(uid, groupInfo.Id, 2) {
		resp.IsDisturb = 1
	}

	return ctx.Success(resp)
}

// UpdateMemberRemark 修改群备注接口
func (c *Group) UpdateMemberRemark(ctx *ichat.Context) error {

	params := &web.GroupRemarkUpdateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	_, err := c.GroupMemberRepo.UpdateWhere(ctx.Ctx(), map[string]any{
		"user_card": params.VisitCard,
	}, "group_id = ? and user_id = ?", params.GroupId, ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness("修改群备注失败！")
	}

	return ctx.Success(nil)
}

func (c *Group) GetInviteFriends(ctx *ichat.Context) error {

	params := &web.GetInviteFriendsRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	items, err := c.ContactService.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if params.GroupId <= 0 {
		return ctx.Success(items)
	}

	mids := c.GroupMemberRepo.GetMemberIds(ctx.Ctx(), int(params.GroupId))
	if len(mids) == 0 {
		return ctx.Success(items)
	}

	data := make([]*model.ContactListItem, 0)
	for i := 0; i < len(items); i++ {
		if !slices.Contains(mids, items[i].Id) {
			data = append(data, items[i])
		}
	}

	return ctx.Success(data)
}

func (c *Group) GroupList(ctx *ichat.Context) error {

	items, err := c.GroupService.List(ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	resp := &web.GroupListResponse{
		Items: make([]*web.GroupListResponse_Item, 0, len(items)),
	}

	for _, item := range items {
		resp.Items = append(resp.Items, &web.GroupListResponse_Item{
			Id:        int32(item.Id),
			GroupName: item.GroupName,
			Avatar:    item.Avatar,
			Profile:   item.Profile,
			Leader:    int32(item.Leader),
			IsDisturb: int32(item.IsDisturb),
			CreatorId: int32(item.CreatorId),
		})
	}

	return ctx.Success(resp)
}

// Members 获取群成员列表
func (c *Group) Members(ctx *ichat.Context) error {

	params := &web.GroupMemberListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	group, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	if group != nil && group.IsDismiss == 1 {
		return ctx.Success([]any{})
	}

	if !c.GroupMemberRepo.IsMember(ctx.Ctx(), int(params.GroupId), ctx.UserId(), false) {
		return ctx.ErrorBusiness("非群成员无权查看成员列表！")
	}

	list := c.GroupMemberRepo.GetMembers(ctx.Ctx(), int(params.GroupId))

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
		})
	}

	return ctx.Success(&web.GroupMemberListResponse{Items: items})
}

// OvertList 公开群列表
func (c *Group) OvertList(ctx *ichat.Context) error {

	params := &web.GroupOvertListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	list, err := c.GroupRepo.SearchOvertList(ctx.Ctx(), &repo.SearchOvertListOpt{
		Name:   params.Name,
		UserId: uid,
		Page:   int(params.Page),
		Size:   20,
	})
	if err != nil {
		return ctx.ErrorBusiness("查询异常！")
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
		return ctx.ErrorBusiness("查询异常！")
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
			Id:        int32(value.Id),
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

// Handover 群主交接
func (c *Group) Handover(ctx *ichat.Context) error {

	params := &web.GroupHandoverRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.GroupMemberRepo.IsMaster(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	if uid == int(params.UserId) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	err := c.GroupMemberService.Handover(ctx.Ctx(), int(params.GroupId), uid, int(params.UserId))
	if err != nil {
		return ctx.ErrorBusiness("转让群主失败！")
	}

	members := make([]model.TalkRecordExtraGroupMembers, 0)
	c.Repo.Db().Table("users").Select("id as user_id", "nickname").Where("id in ?", []int{uid, int(params.UserId)}).Scan(&members)

	extra := model.TalkRecordExtraGroupTransfer{}
	for _, member := range members {
		if member.UserId == uid {
			extra.OldOwnerId = member.UserId
			extra.OldOwnerName = member.Nickname
		} else {
			extra.NewOwnerId = member.UserId
			extra.NewOwnerName = member.Nickname
		}
	}

	_ = c.MessageService.SendSysOther(ctx.Ctx(), &model.TalkRecords{
		MsgType:    entity.ChatMsgSysGroupTransfer,
		TalkType:   model.TalkRecordTalkTypeGroup,
		UserId:     uid,
		ReceiverId: int(params.GroupId),
		Extra:      jsonutil.Encode(extra),
	})

	return ctx.Success(nil)
}

// AssignAdmin 分配管理员
func (c *Group) AssignAdmin(ctx *ichat.Context) error {

	params := &web.GroupAssignAdminRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.GroupMemberRepo.IsMaster(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	leader := 0
	if params.Mode == 1 {
		leader = 1
	}

	err := c.GroupMemberService.SetLeaderStatus(ctx.Ctx(), int(params.GroupId), int(params.UserId), leader)
	if err != nil {
		logger.Errorf("[Group AssignAdmin] 设置管理员信息失败 err :%s", err.Error())
		return ctx.ErrorBusiness("设置管理员信息失败！")
	}

	return ctx.Success(nil)
}

// NoSpeak 禁止发言
func (c *Group) NoSpeak(ctx *ichat.Context) error {

	params := &web.GroupNoSpeakRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	status := 1
	if params.Mode == 2 {
		status = 0
	}

	err := c.GroupMemberService.SetMuteStatus(ctx.Ctx(), int(params.GroupId), int(params.UserId), status)
	if err != nil {
		return ctx.ErrorBusiness("设置群成员禁言状态失败！")
	}

	data := &model.TalkRecords{
		TalkType:   model.TalkRecordTalkTypeGroup,
		UserId:     uid,
		ReceiverId: int(params.GroupId),
	}

	members := make([]model.TalkRecordExtraGroupMembers, 0)
	c.Repo.Db().Table("users").Select("id as user_id", "nickname").Where("id = ?", params.UserId).Scan(&members)

	user, err := c.UsersRepo.FindById(ctx.Ctx(), uid)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if status == 1 {
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

	_ = c.MessageService.SendSysOther(ctx.Ctx(), data)

	return ctx.Success(nil)
}

// Mute 全员禁言
func (c *Group) Mute(ctx *ichat.Context) error {
	params := &web.GroupMuteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	group, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	if group.IsDismiss == 1 {
		return ctx.ErrorBusiness("此群已解散！")
	}

	if !c.GroupMemberRepo.IsLeader(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	data := make(map[string]any)
	if params.Mode == 1 {
		data["is_mute"] = 1
	} else {
		data["is_mute"] = 0
	}

	affected, err := c.GroupRepo.UpdateWhere(ctx.Ctx(), data, "id = ?", params.GroupId)
	if err != nil {
		return ctx.Error("服务器异常，请稍后再试！")
	}

	if affected == 0 {
		return ctx.Success(web.GroupMuteResponse{})
	}

	user, err := c.UsersRepo.FindById(ctx.Ctx(), uid)
	if err != nil {
		return err
	}

	var extra any
	var msgType int
	if params.Mode == 1 {
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

	_ = c.MessageService.SendSysOther(ctx.Ctx(), &model.TalkRecords{
		MsgType:    msgType,
		TalkType:   model.TalkRecordTalkTypeGroup,
		UserId:     uid,
		ReceiverId: int(params.GroupId),
		Extra:      jsonutil.Encode(extra),
	})

	return ctx.Success(web.GroupMuteResponse{})
}

// Overt 公开群
func (c *Group) Overt(ctx *ichat.Context) error {
	params := &web.GroupOvertRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	group, err := c.GroupRepo.FindById(ctx.Ctx(), int(params.GroupId))
	if err != nil {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	if group.IsDismiss == 1 {
		return ctx.ErrorBusiness("此群已解散！")
	}

	if !c.GroupMemberRepo.IsMaster(ctx.Ctx(), int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	data := make(map[string]any)
	if params.Mode == 1 {
		data["is_overt"] = 1
	} else {
		data["is_overt"] = 0
	}

	_, err = c.GroupRepo.UpdateWhere(ctx.Ctx(), data, "id = ?", params.GroupId)
	if err != nil {
		return ctx.Error("服务器异常，请稍后再试！")
	}

	return ctx.Success(web.GroupOvertResponse{})
}
