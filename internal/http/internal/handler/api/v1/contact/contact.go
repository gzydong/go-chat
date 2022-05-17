package contact

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-chat/internal/service/organize"
	"gorm.io/gorm"

	"go-chat/internal/cache"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

type Contact struct {
	service            *service.ContactService
	wsClient           *cache.WsClientSession
	userService        *service.UserService
	talkListService    *service.TalkSessionService
	talkMessageService *service.TalkMessageService
	organizeService    *organize.OrganizeService
}

func NewContactHandler(
	service *service.ContactService,
	wsClient *cache.WsClientSession,
	userService *service.UserService,
	talkListService *service.TalkSessionService,
	talkMessageService *service.TalkMessageService,
	organizeService *organize.OrganizeService,
) *Contact {
	return &Contact{
		service:            service,
		wsClient:           wsClient,
		userService:        userService,
		talkListService:    talkListService,
		talkMessageService: talkMessageService,
		organizeService:    organizeService,
	}
}

// List 联系人列表
func (c *Contact) List(ctx *gin.Context) {
	items, err := c.service.List(ctx, jwtutil.GetUid(ctx))

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

	uid := jwtutil.GetUid(ctx)
	if err := c.service.Delete(ctx, uid, params.FriendId); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	// 解除好友关系后需添加一条聊天记录
	_ = c.talkMessageService.SendSysMessage(ctx, &service.SysTextMessageOpts{
		UserId:     uid,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: params.FriendId,
		Text:       "你与对方已经解除了好友关系！！！",
	})

	// 删除聊天会话
	sid := c.talkListService.Dao().FindBySessionId(uid, params.FriendId, entity.ChatPrivateMode)
	if err := c.talkListService.Delete(ctx, jwtutil.GetUid(ctx), sid); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
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

	response.Success(ctx, entity.H{
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

	if err := c.service.EditRemark(ctx, jwtutil.GetUid(ctx), params.FriendId, params.Remarks); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Detail 联系人详情信息
func (c *Contact) Detail(ctx *gin.Context) {
	params := &request.ContactDetailRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	user, err := c.userService.Dao().FindById(params.UserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.BusinessError(ctx, "用户不存在！")
		} else {
			response.BusinessError(ctx, err)
		}

		return
	}

	resp := entity.H{
		"avatar":          user.Avatar,
		"friend_apply":    0,
		"friend_status":   1, // 朋友关系[0:本人;1:陌生人;2:朋友;]
		"gender":          user.Gender,
		"id":              user.Id,
		"mobile":          user.Mobile,
		"motto":           user.Motto,
		"nickname":        user.Nickname,
		"nickname_remark": "",
	}

	if uid != params.UserId {
		if c.service.Dao().IsFriend(ctx.Request.Context(), uid, params.UserId, false) {
			resp["friend_status"] = 2
			resp["nickname_remark"] = c.service.Dao().GetFriendRemark(ctx.Request.Context(), uid, params.UserId, true)
		} else {
			isOk, _ := c.organizeService.Dao().IsQiyeMember(uid, params.UserId)
			if isOk {
				resp["friend_status"] = 2
			}
		}
	} else {
		resp["friend_status"] = 0
	}

	response.Success(ctx, &resp)
}
