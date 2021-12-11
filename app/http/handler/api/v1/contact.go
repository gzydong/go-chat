package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/strutil"
	"go-chat/app/service"
	"gorm.io/gorm"
	"strconv"
)

type Contact struct {
	service     *service.ContactService
	wsClient    *cache.WsClientSession
	userService *service.UserService
}

func NewContactHandler(
	service *service.ContactService,
	wsClient *cache.WsClientSession,
	userService *service.UserService,
) *Contact {
	return &Contact{
		service:     service,
		wsClient:    wsClient,
		userService: userService,
	}
}

// List 联系人列表
func (c *Contact) List(ctx *gin.Context) {
	items, err := c.service.List(ctx, auth.GetAuthUserID(ctx))

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	for _, item := range items {
		item.IsOnline = strutil.BoolToInt(c.wsClient.IsOnline(ctx, entity.ImChannelDefault, strconv.Itoa(item.Id)))
	}

	response.Success(ctx, items)
}

// Delete 删除联系人
func (c *Contact) Delete(ctx *gin.Context) {
	params := &request.ContactDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.Delete(ctx, auth.GetAuthUserID(ctx), params.FriendId); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Search 查找联系人
func (c *Contact) Search(ctx *gin.Context) {
	params := &request.ContactSearchRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	user, err := c.userService.Dao().FindByMobile(params.Mobile)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.BusinessError(ctx, "用户不存在！")
		} else {
			response.BusinessError(ctx, err)
		}

		return
	}

	response.Success(ctx, &gin.H{
		"id":       user.Id,
		"mobile":   user.Mobile,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"gender":   user.Gender,
		"motto":    user.Motto,
	})
}

// EditRemark 编辑联系人备注
func (c *Contact) EditRemark(ctx *gin.Context) {
	params := &request.ContactEditRemarkRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.EditRemark(ctx, auth.GetAuthUserID(ctx), params.FriendId, params.Remarks); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{})
}

// Detail 联系人详情信息
func (c *Contact) Detail(ctx *gin.Context) {
	params := &request.ContactDetailRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	user, err := c.userService.Dao().FindById(params.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.BusinessError(ctx, "用户不存在！")
		} else {
			response.BusinessError(ctx, err)
		}

		return
	}

	resp := gin.H{
		"avatar":          user.Avatar,
		"friend_apply":    0,
		"friend_status":   0,
		"gender":          user.Gender,
		"id":              user.Id,
		"mobile":          user.Mobile,
		"motto":           user.Motto,
		"nickname":        user.Nickname,
		"nickname_remark": "",
	}

	response.Success(ctx, &resp)
}
