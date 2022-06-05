package contact

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/timeutil"

	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service"
)

type ContactApply struct {
	service            *service.ContactApplyService
	userService        *service.UserService
	talkMessageService *service.TalkMessageService
	contactService     *service.ContactService
}

func NewContactsApplyHandler(
	service *service.ContactApplyService,
	userService *service.UserService,
	talkMessageService *service.TalkMessageService,
	contactService *service.ContactService,
) *ContactApply {
	return &ContactApply{service: service, userService: userService, talkMessageService: talkMessageService, contactService: contactService}
}

// ApplyUnreadNum 获取好友申请未读数
func (c *ContactApply) ApplyUnreadNum(ctx *gin.Context) {
	response.Success(ctx, entity.H{
		"unread_num": c.service.GetApplyUnreadNum(ctx.Request.Context(), jwtutil.GetUid(ctx)),
	})
}

// Create 创建联系人申请
func (c *ContactApply) Create(ctx *gin.Context) {
	params := &request.ContactApplyCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)
	if !c.contactService.Dao().IsFriend(ctx, uid, params.FriendId, false) {
		response.Success(ctx, nil)
	}

	if err := c.service.Create(ctx, &service.ContactApplyCreateOpts{
		UserId:   jwtutil.GetUid(ctx),
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

	uid := jwtutil.GetUid(ctx)
	applyInfo, err := c.service.Accept(ctx, &service.ContactApplyAcceptOpts{
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
		UserId:  uid,
	})

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	_ = c.talkMessageService.SendSysMessage(ctx, &service.SysTextMessageOpts{
		UserId:     applyInfo.UserId,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: applyInfo.FriendId,
		Text:       "你们已成为好友，可以开始聊天咯！",
	})

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
		UserId:  jwtutil.GetUid(ctx),
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
	list, err := c.service.List(ctx, jwtutil.GetUid(ctx), 1, 1000)
	if err != nil {
		response.SystemError(ctx, err)
		return
	}

	items := make([]*entity.H, 0)
	for _, item := range list {
		items = append(items, &entity.H{
			"id":         item.Id,
			"user_id":    item.UserId,
			"friend_id":  item.FriendId,
			"remark":     item.Remark,
			"nickname":   item.Nickname,
			"avatar":     item.Avatar,
			"created_at": timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	c.service.ClearApplyUnreadNum(ctx, jwtutil.GetUid(ctx))

	response.SuccessPaginate(ctx, items, 1, 1000, len(items))
}
