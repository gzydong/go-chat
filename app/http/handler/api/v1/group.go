package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service"
)

type Group struct {
	service         *service.GroupService
	memberService   *service.GroupMemberService
	talkListService *service.TalkSessionService
	userService     *service.UserService
	redisLock       *cache.RedisLock
	contactService  *service.ContactService
}

func NewGroupHandler(
	service *service.GroupService,
	memberService *service.GroupMemberService,
	talkListService *service.TalkSessionService,
	redisLock *cache.RedisLock,
	contactService *service.ContactService,
	userService *service.UserService,
) *Group {
	return &Group{
		service:         service,
		memberService:   memberService,
		talkListService: talkListService,
		redisLock:       redisLock,
		contactService:  contactService,
		userService:     userService,
	}
}

// Create 创建群聊分组
func (c *Group) Create(ctx *gin.Context) {
	params := &request.GroupCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	gid, err := c.service.Create(ctx, params)
	if err != nil {
		response.BusinessError(ctx, "创建群聊失败，请稍后再试！")
		return
	}

	response.Success(ctx, gin.H{
		"group_id": gid,
	})
}

// Dismiss 解散群组
func (c *Group) Dismiss(ctx *gin.Context) {
	params := &request.GroupDismissRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Dismiss(params.GroupId, auth.GetAuthUserID(ctx)); err != nil {
		response.BusinessError(ctx, "群组解散失败！")
		return
	}

	response.Success(ctx, nil)
}

// Invite 邀请好友加入群聊
func (c *Group) Invite(ctx *gin.Context) {
	params := &request.GroupInviteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	key := fmt.Sprintf("group-join:%d", params.GroupId)
	if !c.redisLock.Lock(ctx, key, 20) {
		response.BusinessError(ctx, "网络异常，请稍后再试！")
		return
	}

	// 释放锁
	defer c.redisLock.Release(ctx, key)

	uid := auth.GetAuthUserID(ctx)
	uids := slice.UniqueInt(slice.ParseIds(params.Ids))

	if len(uids) == 0 {
		response.BusinessError(ctx, "邀请好友列表不能为空！")
		return
	}

	if !c.memberService.Dao().IsMember(params.GroupId, uid, true) {
		response.BusinessError(ctx, "非群组成员，无权邀请好友！")
		return
	}

	if err := c.service.InviteUsers(ctx, params.GroupId, uid, uids); err != nil {
		response.BusinessError(ctx, "邀请好友加入群聊失败！")
		return
	}

	response.Success(ctx, nil)
}

// SignOut 退出群聊
func (c *Group) SignOut(ctx *gin.Context) {
	params := &request.GroupSecedeRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Secede(params.GroupId, auth.GetAuthUserID(ctx)); err != nil {
		response.BusinessError(ctx, "退出群组失败！")
		return
	}

	response.Success(ctx, nil)
}

// Setting 群设置接口（预留）
func (c *Group) Setting(ctx *gin.Context) {
	params := &request.GroupSettingRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !c.memberService.Dao().IsLeader(params.GroupId, auth.GetAuthUserID(ctx)) {
		response.BusinessError(ctx, "无权限操作")
		return
	}

	if err := c.service.Update(ctx.Request.Context(), params); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// RemoveMembers 移除指定成员(群组&管理员权限)
func (c *Group) RemoveMembers(ctx *gin.Context) {
	params := &request.GroupRemoveMembersRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !c.memberService.Dao().IsLeader(params.GroupId, auth.GetAuthUserID(ctx)) {
		response.BusinessError(ctx, "无权限操作")
		return
	}

	// todo 移除指定成员

	response.Success(ctx, nil)
}

// Detail 获取群组信息
func (c *Group) Detail(ctx *gin.Context) {
	params := &request.GroupCommonRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	groupInfo, err := c.service.Dao().FindById(params.GroupId)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	if groupInfo.Id == 0 {
		response.BusinessError(ctx, "数据不存在")
		return
	}

	info := gin.H{}
	info["group_id"] = groupInfo.Id
	info["group_name"] = groupInfo.Name
	info["profile"] = groupInfo.Profile
	info["avatar"] = groupInfo.Avatar
	info["created_at"] = timeutil.FormatDatetime(groupInfo.CreatedAt)
	info["is_manager"] = uid == groupInfo.CreatorId
	info["manager_nickname"] = ""
	info["visit_card"] = c.memberService.Dao().GetMemberRemark(params.GroupId, uid)
	info["is_disturb"] = 0
	info["notice"] = []gin.H{}

	if c.talkListService.Dao().IsDisturb(uid, groupInfo.Id, 2) {
		info["is_disturb"] = 1
	}

	if userInfo, err := c.userService.Dao().FindById(uid); err == nil {
		info["manager_nickname"] = userInfo.Nickname
	}

	response.Success(ctx, info)
}

// EditRemark 修改群备注接口
func (c *Group) EditRemark(ctx *gin.Context) {
	params := &request.GroupEditRemarkRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.memberService.UpdateMemberCard(params.GroupId, auth.GetAuthUserID(ctx), params.VisitCard); err != nil {
		response.BusinessError(ctx, "修改群备注失败！")
		return
	}

	response.Success(ctx, nil)
}

func (c *Group) GetInviteFriends(ctx *gin.Context) {
	params := &request.GetInviteFriendsRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	items, err := c.contactService.List(ctx, auth.GetAuthUserID(ctx))
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	if params.GroupId <= 0 {
		response.Success(ctx, items)
		return
	}

	mids := c.memberService.Dao().GetMemberIds(params.GroupId)
	if len(mids) == 0 {
		response.Success(ctx, items)
		return
	}

	data := make([]*model.ContactListItem, 0)
	for i := 0; i < len(items); i++ {
		if !slice.InInt(items[i].Id, mids) {
			data = append(data, items[i])
		}
	}

	response.Success(ctx, data)
}

func (c *Group) GetGroups(ctx *gin.Context) {
	items, err := c.service.UserGroupList(auth.GetAuthUserID(ctx))
	if err != nil {
		response.BusinessError(ctx, items)
		return
	}

	response.Success(ctx, items)
}

// GetGroupMembers 获取群成员列表
func (c *Group) GetGroupMembers(ctx *gin.Context) {
	params := &request.GroupCommonRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !c.memberService.Dao().IsMember(params.GroupId, auth.GetAuthUserID(ctx), false) {
		response.BusinessError(ctx, "非群成员无权查看成员列表！")
		return
	}

	response.Success(ctx, c.memberService.Dao().GetMembers(params.GroupId))
}
