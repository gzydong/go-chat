package group

import (
	"fmt"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type Group struct {
	service            *service.GroupService
	memberService      *service.GroupMemberService
	talkListService    *service.TalkSessionService
	userService        *service.UserService
	redisLock          *cache.RedisLock
	contactService     *service.ContactService
	groupNoticeService *service.GroupNoticeService
	messageService     *service.TalkMessageService
}

func NewGroup(service *service.GroupService, memberService *service.GroupMemberService, talkListService *service.TalkSessionService, userService *service.UserService, redisLock *cache.RedisLock, contactService *service.ContactService, groupNoticeService *service.GroupNoticeService, messageService *service.TalkMessageService) *Group {
	return &Group{service: service, memberService: memberService, talkListService: talkListService, userService: userService, redisLock: redisLock, contactService: contactService, groupNoticeService: groupNoticeService, messageService: messageService}
}

// Create 创建群聊分组
func (c *Group) Create(ctx *ichat.Context) error {

	params := &web.GroupCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	gid, err := c.service.Create(ctx.Ctx(), &service.CreateGroupOpt{
		UserId:    ctx.UserId(),
		Name:      params.Name,
		Avatar:    params.Avatar,
		MemberIds: sliceutil.ParseIds(params.GetIds()),
	})
	if err != nil {
		return ctx.ErrorBusiness("创建群聊失败，请稍后再试！")
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
	if !c.memberService.Dao().IsMaster(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限解散群组！")
	}

	if err := c.service.Dismiss(ctx.Ctx(), int(params.GroupId), ctx.UserId()); err != nil {
		return ctx.ErrorBusiness("群组解散失败！")
	}

	_ = c.messageService.SendSysMessage(ctx.Context, &service.SysTextMessageOpt{
		UserId:     uid,
		TalkType:   entity.ChatGroupMode,
		ReceiverId: int(params.GroupId),
		Text:       "群组已被群主或管理员解散！",
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
	if !c.redisLock.Lock(ctx.Context, key, 20) {
		return ctx.ErrorBusiness("网络异常，请稍后再试！")
	}

	defer c.redisLock.UnLock(ctx.Context, key)

	uid := ctx.UserId()
	uids := sliceutil.Unique(sliceutil.ParseIds(params.Ids))

	if len(uids) == 0 {
		return ctx.ErrorBusiness("邀请好友列表不能为空！")
	}

	if !c.memberService.Dao().IsMember(int(params.GroupId), uid, true) {
		return ctx.ErrorBusiness("非群组成员，无权邀请好友！")
	}

	if err := c.service.InviteMembers(ctx.Context, &service.InviteGroupMembersOpt{
		UserId:    uid,
		GroupId:   int(params.GroupId),
		MemberIds: uids,
	}); err != nil {
		return ctx.ErrorBusiness("邀请好友加入群聊失败！")
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
	if err := c.service.Secede(ctx.Ctx(), int(params.GroupId), uid); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	// 删除聊天会话
	sid := c.talkListService.Dao().FindBySessionId(uid, int(params.GroupId), entity.ChatGroupMode)
	_ = c.talkListService.Delete(ctx.Context, ctx.UserId(), sid)

	return ctx.Success(nil)
}

// Setting 群设置接口（预留）
func (c *Group) Setting(ctx *ichat.Context) error {

	params := &web.GroupSettingRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.memberService.Dao().IsLeader(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("无权限操作")
	}

	if err := c.service.Update(ctx.Ctx(), &service.UpdateGroupOpt{
		GroupId: int(params.GroupId),
		Name:    params.GroupName,
		Avatar:  params.Avatar,
		Profile: params.Profile,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	_ = c.messageService.SendSysMessage(ctx.Context, &service.SysTextMessageOpt{
		UserId:     uid,
		TalkType:   entity.ChatGroupMode,
		ReceiverId: int(params.GroupId),
		Text:       "群主或管理员修改了群信息！",
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

	if !c.memberService.Dao().IsLeader(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("无权限操作")
	}

	err := c.service.RemoveMembers(ctx.Ctx(), &service.RemoveMembersOpt{
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

	groupInfo, err := c.service.Dao().FindById(ctx.Ctx(), int(params.GroupId))
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
		VisitCard: c.memberService.Dao().GetMemberRemark(int(params.GroupId), uid),
	}

	if c.talkListService.Dao().IsDisturb(uid, groupInfo.Id, 2) {
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

	_, err := c.memberService.Dao().Updates(ctx.Ctx(), map[string]interface{}{
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

	items, err := c.contactService.List(ctx.Context, ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if params.GroupId <= 0 {
		return ctx.Success(items)
	}

	mids := c.memberService.Dao().GetMemberIds(int(params.GroupId))
	if len(mids) == 0 {
		return ctx.Success(items)
	}

	data := make([]*model.ContactListItem, 0)
	for i := 0; i < len(items); i++ {
		if !sliceutil.Include(items[i].Id, mids) {
			data = append(data, items[i])
		}
	}

	return ctx.Success(data)
}

func (c *Group) GroupList(ctx *ichat.Context) error {

	items, err := c.service.List(ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	resp := &web.GroupListResponse{
		Rows: make([]*web.GroupListResponse_Item, 0, len(items)),
	}

	for _, item := range items {
		resp.Rows = append(resp.Rows, &web.GroupListResponse_Item{
			Id:        int32(item.Id),
			GroupName: item.GroupName,
			Avatar:    item.Avatar,
			Profile:   item.Profile,
			Leader:    int32(item.Leader),
			IsDisturb: int32(item.IsDisturb),
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

	if !c.memberService.Dao().IsMember(int(params.GroupId), ctx.UserId(), false) {
		return ctx.ErrorBusiness("非群成员无权查看成员列表！")
	}

	return ctx.Success(c.memberService.Dao().GetMembers(int(params.GroupId)))
}

// OvertList 公开群列表
func (c *Group) OvertList(ctx *ichat.Context) error {

	params := &web.GroupOvertListRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	list, err := c.service.Dao().SearchOvertList(ctx.Context, params.Name, int(params.Page), 20)
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

	count, err := c.memberService.Dao().CountGroupMemberNum(ids)
	if err != nil {
		return ctx.ErrorBusiness("查询异常！")
	}

	countMap := make(map[int]int)
	for _, member := range count {
		countMap[member.GroupId] = member.Count
	}

	checks, err := c.memberService.Dao().CheckUserGroup(ids, ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness("查询异常！")
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
			IsMember:  sliceutil.Include(value.Id, checks),
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
	if !c.memberService.Dao().IsMaster(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	if uid == int(params.UserId) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	err := c.memberService.Handover(ctx.Ctx(), int(params.GroupId), uid, int(params.UserId))
	if err != nil {
		return ctx.ErrorBusiness("转让群主失败！")
	}

	return ctx.Success(entity.H{})
}

// AssignAdmin 分配管理员
func (c *Group) AssignAdmin(ctx *ichat.Context) error {

	params := &web.GroupAssignAdminRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.memberService.Dao().IsMaster(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	leader := 0
	if params.Mode == 1 {
		leader = 1
	}

	err := c.memberService.UpdateLeaderStatus(int(params.GroupId), int(params.UserId), leader)
	if err != nil {
		logger.Error("[Group AssignAdmin] 设置管理员信息失败 err :", err.Error())
		return ctx.ErrorBusiness("设置管理员信息失败！")
	}

	return ctx.Success(entity.H{})
}

// NoSpeak 禁止发言
func (c *Group) NoSpeak(ctx *ichat.Context) error {

	params := &web.GroupNoSpeakRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.memberService.Dao().IsLeader(int(params.GroupId), uid) {
		return ctx.ErrorBusiness("暂无权限！")
	}

	status := 1
	if params.Mode == 2 {
		status = 0
	}

	err := c.memberService.UpdateMuteStatus(int(params.GroupId), int(params.UserId), status)
	if err != nil {
		return ctx.ErrorBusiness("设置群成员禁言状态失败！")
	}

	return ctx.Success(entity.H{})
}
