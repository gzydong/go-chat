package contact

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ginutil"
	"go-chat/internal/pkg/timeutil"

	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/service"
)

type Apply struct {
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
) *Apply {
	return &Apply{service: service, userService: userService, talkMessageService: talkMessageService, contactService: contactService}
}

// ApplyUnreadNum 获取好友申请未读数
func (c *Apply) ApplyUnreadNum(ctx *gin.Context) error {
	return ginutil.Success(ctx, entity.H{
		"unread_num": c.service.GetApplyUnreadNum(ctx.Request.Context(), jwtutil.GetUid(ctx)),
	})
}

// Create 创建联系人申请
func (c *Apply) Create(ctx *gin.Context) error {
	params := &web.ContactApplyCreateRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)
	if !c.contactService.Dao().IsFriend(ctx, uid, params.FriendId, false) {
		return ginutil.Success(ctx, nil)
	}

	if err := c.service.Create(ctx, &service.ContactApplyCreateOpts{
		UserId:   jwtutil.GetUid(ctx),
		Remarks:  params.Remarks,
		FriendId: params.FriendId,
	}); err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, nil)
}

// Accept 同意联系人添加申请
func (c *Apply) Accept(ctx *gin.Context) error {
	params := &web.ContactApplyAcceptRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)
	applyInfo, err := c.service.Accept(ctx, &service.ContactApplyAcceptOpts{
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
		UserId:  uid,
	})

	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	_ = c.talkMessageService.SendSysMessage(ctx, &service.SysTextMessageOpts{
		UserId:     applyInfo.UserId,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: applyInfo.FriendId,
		Text:       "你们已成为好友，可以开始聊天咯！",
	})

	return ginutil.Success(ctx, nil)
}

// Decline 拒绝联系人添加申请
func (c *Apply) Decline(ctx *gin.Context) error {
	params := &web.ContactApplyDeclineRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	if err := c.service.Decline(ctx, &service.ContactApplyDeclineOpts{
		UserId:  jwtutil.GetUid(ctx),
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
	}); err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, nil)
}

// List 获取联系人申请列表
func (c *Apply) List(ctx *gin.Context) error {
	list, err := c.service.List(ctx, jwtutil.GetUid(ctx), 1, 1000)
	if err != nil {
		return ginutil.SystemError(ctx, err)
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

	return ginutil.SuccessPaginate(ctx, items, 1, 1000, len(items))
}
