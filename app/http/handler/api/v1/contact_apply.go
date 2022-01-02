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
	response.Success(ctx, gin.H{
		"unread_num": 0,
	})
}

// Create 创建联系人申请
func (c *ContactApply) Create(ctx *gin.Context) {
	params := &request.ContactApplyCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Create(ctx, &service.ContactApplyCreateOpts{
		UserId:   auth.GetAuthUserID(ctx),
		Remarks:  params.Remarks,
		FriendId: params.FriendId,
	}); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Accept 同意联系人添加申请
func (c *ContactApply) Accept(ctx *gin.Context) {
	params := &request.ContactApplyAcceptRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Accept(ctx, &service.ContactApplyAcceptOpts{
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
	}); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Decline 拒绝联系人添加申请
func (c *ContactApply) Decline(ctx *gin.Context) {
	params := &request.ContactApplyDeclineRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Decline(ctx, &service.ContactApplyDeclineOpts{
		UserId:  auth.GetAuthUserID(ctx),
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
	}); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// List 获取联系人申请列表
func (c *ContactApply) List(ctx *gin.Context) {
	items, err := c.service.List(ctx, auth.GetAuthUserID(ctx), 1, 1000)
	if err != nil {
		response.SystemError(ctx, err)
		return
	}

	response.SuccessPaginate(ctx, items, 1, 1000, len(items))
}
