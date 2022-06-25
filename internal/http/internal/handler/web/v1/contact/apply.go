package contact

import (
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
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
func (c *Apply) ApplyUnreadNum(ctx *ichat.Context) error {
	return ctx.Success(entity.H{
		"unread_num": c.service.GetApplyUnreadNum(ctx.Context.Request.Context(), jwtutil.GetUid(ctx.Context)),
	})
}

// Create 创建联系人申请
func (c *Apply) Create(ctx *ichat.Context) error {
	params := &web.ContactApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := jwtutil.GetUid(ctx.Context)
	if !c.contactService.Dao().IsFriend(ctx.Context, uid, params.FriendId, false) {
		return ctx.Success(nil)
	}

	if err := c.service.Create(ctx.Context, &service.ContactApplyCreateOpts{
		UserId:   jwtutil.GetUid(ctx.Context),
		Remarks:  params.Remarks,
		FriendId: params.FriendId,
	}); err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(nil)
}

// Accept 同意联系人添加申请
func (c *Apply) Accept(ctx *ichat.Context) error {

	params := &web.ContactApplyAcceptRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := jwtutil.GetUid(ctx.Context)
	applyInfo, err := c.service.Accept(ctx.Context, &service.ContactApplyAcceptOpts{
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
		UserId:  uid,
	})

	if err != nil {
		return ctx.BusinessError(err)
	}

	_ = c.talkMessageService.SendSysMessage(ctx.Context, &service.SysTextMessageOpts{
		UserId:     applyInfo.UserId,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: applyInfo.FriendId,
		Text:       "你们已成为好友，可以开始聊天咯！",
	})

	return ctx.Success(nil)
}

// Decline 拒绝联系人添加申请
func (c *Apply) Decline(ctx *ichat.Context) error {

	params := &web.ContactApplyDeclineRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.Decline(ctx.Context, &service.ContactApplyDeclineOpts{
		UserId:  jwtutil.GetUid(ctx.Context),
		Remarks: params.Remarks,
		ApplyId: params.ApplyId,
	}); err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(nil)
}

// List 获取联系人申请列表
func (c *Apply) List(ctx *ichat.Context) error {

	list, err := c.service.List(ctx.Context, jwtutil.GetUid(ctx.Context), 1, 1000)
	if err != nil {
		return ctx.Error(err.Error())
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

	c.service.ClearApplyUnreadNum(ctx.Context, jwtutil.GetUid(ctx.Context))

	return ctx.Paginate(items, 1, 1000, len(items))
}
