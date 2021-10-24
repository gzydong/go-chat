package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/service"
)

type Group struct {
	service       *service.GroupService
	memberService *service.GroupMemberService
}

func NewGroupHandler(service *service.GroupService, memberService *service.GroupMemberService) *Group {
	return &Group{
		service:       service,
		memberService: memberService,
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

func (c *Group) Invite(ctx *gin.Context) {

}

func (c *Group) Secede(ctx *gin.Context) {
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

func (c *Group) Setting(ctx *gin.Context) {

}

func (c *Group) RemoveMembers(ctx *gin.Context) {

}

func (c *Group) Detail(ctx *gin.Context) {
	params := &request.GroupCommonRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	info := make(map[string]interface{})

	info["group_id"] = params.GroupId
	info["group_name"] = ""
	info["profile"] = ""
	info["avatar"] = ""
	info["created_at"] = ""
	info["is_manager"] = ""
	info["manager_nickname"] = ""
	info["visit_card"] = c.memberService.GetMemberRemarks(params.GroupId, auth.GetAuthUserID(ctx))
	info["is_disturb"] = ""
	info["notice"] = ""

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
