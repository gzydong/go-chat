package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/dao"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service"
)

type Group struct {
	service         *service.GroupService
	memberService   *service.GroupMemberService
	talkListService *service.TalkListService
	userRepo        *dao.UserDao
	redisLock       *cache.RedisLock
}

func NewGroupHandler(
	service *service.GroupService,
	memberService *service.GroupMemberService,
	talkListService *service.TalkListService,
	userRepo *dao.UserDao,
	redisLock *cache.RedisLock,
) *Group {
	return &Group{
		service:         service,
		memberService:   memberService,
		talkListService: talkListService,
		userRepo:        userRepo,
		redisLock:       redisLock,
	}
}

// Create 创建群聊分组
func (c *Group) Create(ctx *gin.Context) {
	params := &request.GroupCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Create(ctx, params); err != nil {
		response.BusinessError(ctx, "创建群聊失败，请稍后再试！")
		return
	}

	response.Success(ctx, gin.H{})
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

	response.Success(ctx, gin.H{})
}

// Invite 邀请好友加入群聊
func (c *Group) Invite(ctx *gin.Context) {
	params := &request.GroupInviteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	keyLock := fmt.Sprintf("group-invite:%d", params.GroupId)

	if !c.redisLock.Lock(ctx, keyLock, 20) {
		response.BusinessError(ctx, "网络异常，请稍后再试！")
		return
	}

	// 释放锁
	defer c.redisLock.Release(ctx, keyLock)

	uid := auth.GetAuthUserID(ctx)
	uids := slice.UniqueInt(slice.ParseIds(params.Ids))

	if len(uids) == 0 {
		response.BusinessError(ctx, "邀请好友列表不能为空！")
		return
	}

	if !c.memberService.IsMember(params.GroupId, uid) {
		response.BusinessError(ctx, "非群组成员，无权邀请好友！")
		return
	}

	if err := c.service.InviteUsers(params.GroupId, uid, uids); err != nil {
		response.BusinessError(ctx, "邀请好友加入群聊失败！")
		return
	}

	response.Success(ctx, gin.H{}, "邀请成功！")
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

	response.Success(ctx, gin.H{})
}

// Setting 群设置接口（预留）
func (c *Group) Setting(ctx *gin.Context) {

}

// RemoveMembers 移除指定成员(群组&管理员权限)
func (c *Group) RemoveMembers(ctx *gin.Context) {
	params := &request.GroupRemoveMembersRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}
}

// Detail 获取群组信息
func (c *Group) Detail(ctx *gin.Context) {
	params := &request.GroupCommonRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	groupInfo, err := c.service.FindById(params.GroupId)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	if groupInfo.ID == 0 {
		response.BusinessError(ctx, "数据不存在")
		return
	}

	info := gin.H{}
	info["group_id"] = groupInfo.ID
	info["group_name"] = groupInfo.GroupName
	info["profile"] = groupInfo.Profile
	info["avatar"] = groupInfo.Avatar
	info["created_at"] = timeutil.FormatDatetime(groupInfo.CreatedAt)
	info["is_manager"] = uid == groupInfo.CreatorId
	info["manager_nickname"] = ""
	info["visit_card"] = c.memberService.GetMemberRemarks(params.GroupId, uid)
	info["is_disturb"] = 0
	info["notice"] = []gin.H{}

	if c.talkListService.IsDisturb(uid, groupInfo.ID, 2) {
		info["is_disturb"] = 1
	}

	if userInfo, err := c.userRepo.FindById(uid); err == nil {
		info["manager_nickname"] = userInfo.Nickname
	}

	response.Success(ctx, info)
}

// EditGroupCard 修改群备注接口
func (c *Group) EditGroupRemarks(ctx *gin.Context) {
	params := &request.GroupEditCardRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.UpdateMemberCard(params.GroupId, auth.GetAuthUserID(ctx), params.VisitCard); err != nil {
		response.BusinessError(ctx, "修改群备注失败！")
		return
	}

	response.Success(ctx, gin.H{})
}

func (c *Group) GetInviteFriends(ctx *gin.Context) {

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

	if !c.memberService.IsMember(params.GroupId, auth.GetAuthUserID(ctx)) {
		response.BusinessError(ctx, "非群成员无权查看成员列表！")
		return
	}

	items := c.memberService.GetGroupMembers(params.GroupId)

	response.Success(ctx, items)
}
