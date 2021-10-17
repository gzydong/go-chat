package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/service"
)

type Group struct {
	service *service.GroupService
}

func NewGroupHandler(groupService *service.GroupService) *Group {
	return &Group{
		service: groupService,
	}
}

// Create 创建群聊分组
func (c *Group) Create(ctx *gin.Context) {
	params := &request.GroupCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Create(ctx, params)
	if err != nil {
		response.BusinessError(ctx, "创建群聊失败，请稍后再试！")
		return
	}

	response.Success(ctx, gin.H{})
}

func (c *Group) Dismiss(ctx *gin.Context) {
	params := &request.GroupDismissRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Dismiss(params.GroupId, helper.GetAuthUserID(ctx))
	if err != nil {
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

	err := c.service.Secede(params.GroupId, helper.GetAuthUserID(ctx))
	if err != nil {
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

}

// EditGroupCard 修改群备注接口
func (c *Group) EditGroupCard(ctx *gin.Context) {
	params := &request.GroupEditCardRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.UpdateMemberCard(params.GroupId, helper.GetAuthUserID(ctx), params.VisitCard)
	if err != nil {
		response.BusinessError(ctx, "修改群备注失败！")
		return
	}

	response.Success(ctx, gin.H{})
}

func (c *Group) GetInviteFriends(ctx *gin.Context) {

}

func (c *Group) GetGroups(ctx *gin.Context) {
	items, err := c.service.UserGroupList(helper.GetAuthUserID(ctx))
	if err != nil {
		response.BusinessError(ctx, err.Error())
		return
	}

	response.Success(ctx, items)
}

func (c *Group) GetGroupMembers(ctx *gin.Context) {

}

func (c *Group) GetGroupNotice(ctx *gin.Context) {

}

func (c *Group) EditNotice(ctx *gin.Context) {

}

func (c *Group) DeleteNotice(ctx *gin.Context) {

}
