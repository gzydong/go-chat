package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/service"
)

type ContactApply struct {
	service     *service.ContactApplyService
	userService *service.UserService
}

func NewContactsApplyHandler(
	service *service.ContactApplyService,
	userService *service.UserService,
) *ContactApply {
	return &ContactApply{service: service, userService: userService}
}

// ApplyUnreadNum 获取好友申请未读数
func (c *ContactApply) ApplyUnreadNum(ctx *gin.Context) {

}

// Create 创建联系人申请
func (c *ContactApply) Create(ctx *gin.Context) {
	params := &request.ContactApplyCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Create(ctx, auth.GetAuthUserID(ctx), params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Accept 同意联系人添加申请
func (c *ContactApply) Accept(ctx *gin.Context) {
	params := &request.ContactApplyAcceptRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Accept(ctx, auth.GetAuthUserID(ctx), params); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	// todo 推送消息

	response.Success(ctx, gin.H{})
}

// Decline 拒绝联系人添加申请
func (c *ContactApply) Decline(ctx *gin.Context) {

}

// List 获取联系人申请列表
func (c *ContactApply) List(ctx *gin.Context) {
	items, err := c.service.List(ctx, auth.GetAuthUserID(ctx), 1, 1000)
	if err != nil {
		response.SystemError(ctx, err)
	}

	response.SuccessPaginate(ctx, items, 1, 1000, len(items))
}
