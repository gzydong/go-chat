package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/service"
)

type Apply struct {
	service            *service.ContactApplyService
	userService        *service.UserService
	talkMessageService *service.TalkMessageService
	contactService     *service.ContactService
}

func NewApply(service *service.ContactApplyService, userService *service.UserService, talkMessageService *service.TalkMessageService, contactService *service.ContactService) *Apply {
	return &Apply{service: service, userService: userService, talkMessageService: talkMessageService, contactService: contactService}
}

// ApplyUnreadNum 获取好友申请未读数
func (c *Apply) ApplyUnreadNum(ctx *ichat.Context) error {
	return ctx.Success(entity.H{
		"unread_num": c.service.GetApplyUnreadNum(ctx.Ctx(), ctx.UserId()),
	})
}

// Create 创建联系人申请
func (c *Apply) Create(ctx *ichat.Context) error {

	params := &web.ContactApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if !c.contactService.Dao().IsFriend(ctx.Context, uid, int(params.FriendId), false) {
		return ctx.Success(nil)
	}

	if err := c.service.Create(ctx.Context, &service.ContactApplyCreateOpts{
		UserId:   ctx.UserId(),
		Remarks:  params.Remark,
		FriendId: int(params.FriendId),
	}); err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(&web.ContactApplyCreateResponse{})
}

// Accept 同意联系人添加申请
func (c *Apply) Accept(ctx *ichat.Context) error {

	params := &web.ContactApplyAcceptRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	applyInfo, err := c.service.Accept(ctx.Context, &service.ContactApplyAcceptOpts{
		Remarks: params.Remark,
		ApplyId: int(params.ApplyId),
		UserId:  uid,
	})

	if err != nil {
		return ctx.BusinessError(err)
	}

	_ = c.talkMessageService.SendSysMessage(ctx.Context, &service.SysTextMessageOpt{
		UserId:     applyInfo.UserId,
		TalkType:   entity.ChatPrivateMode,
		ReceiverId: applyInfo.FriendId,
		Text:       "你们已成为好友，可以开始聊天咯！",
	})

	return ctx.Success(&web.ContactApplyAcceptResponse{})
}

// Decline 拒绝联系人添加申请
func (c *Apply) Decline(ctx *ichat.Context) error {

	params := &web.ContactApplyDeclineRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.Decline(ctx.Context, &service.ContactApplyDeclineOpts{
		UserId:  ctx.UserId(),
		Remarks: params.Remark,
		ApplyId: int(params.ApplyId),
	}); err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(&web.ContactApplyDeclineResponse{})
}

// List 获取联系人申请列表
func (c *Apply) List(ctx *ichat.Context) error {

	list, err := c.service.List(ctx.Context, ctx.UserId(), 1, 1000)
	if err != nil {
		return ctx.Error(err.Error())
	}

	items := make([]*web.ContactApplyListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ContactApplyListResponse_Item{
			Id:        int32(item.Id),
			UserId:    int32(item.UserId),
			FriendId:  int32(item.FriendId),
			Remark:    item.Remark,
			Nickname:  item.Nickname,
			Avatar:    item.Avatar,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	c.service.ClearApplyUnreadNum(ctx.Context, ctx.UserId())

	return ctx.Success(&web.ContactApplyListResponse{
		Items: items,
	})
}
