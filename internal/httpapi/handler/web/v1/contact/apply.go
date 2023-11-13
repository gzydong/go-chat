package contact

import (
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Apply struct {
	ContactRepo         *repo.Contact
	ContactApplyService service.IContactApplyService
	UserService         service.IUserService
	ContactService      service.IContactService
	MessageService      service.IMessageService
}

// ApplyUnreadNum 获取好友申请未读数
func (c *Apply) ApplyUnreadNum(ctx *ichat.Context) error {
	return ctx.Success(map[string]any{
		"unread_num": c.ContactApplyService.GetApplyUnreadNum(ctx.Ctx(), ctx.UserId()),
	})
}

// Create 创建联系人申请
func (c *Apply) Create(ctx *ichat.Context) error {

	params := &web.ContactApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if c.ContactRepo.IsFriend(ctx.Ctx(), uid, int(params.FriendId), false) {
		return ctx.Success(nil)
	}

	if err := c.ContactApplyService.Create(ctx.Ctx(), &service.ContactApplyCreateOpt{
		UserId:   ctx.UserId(),
		Remarks:  params.Remark,
		FriendId: int(params.FriendId),
	}); err != nil {
		return ctx.ErrorBusiness(err)
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
	applyInfo, err := c.ContactApplyService.Accept(ctx.Ctx(), &service.ContactApplyAcceptOpt{
		Remarks: params.Remark,
		ApplyId: int(params.ApplyId),
		UserId:  uid,
	})

	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	err = c.MessageService.SendSystemText(ctx.Ctx(), applyInfo.UserId, &message.TextMessageRequest{
		Content: "你们已成为好友，可以开始聊天咯！",
		Receiver: &message.MessageReceiver{
			TalkType:   entity.ChatPrivateMode,
			ReceiverId: int32(applyInfo.FriendId),
		},
	})

	if err != nil {
		fmt.Println("Apply Accept Err", err.Error())
	}

	return ctx.Success(&web.ContactApplyAcceptResponse{})
}

// Decline 拒绝联系人添加申请
func (c *Apply) Decline(ctx *ichat.Context) error {

	params := &web.ContactApplyDeclineRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ContactApplyService.Decline(ctx.Ctx(), &service.ContactApplyDeclineOpt{
		UserId:  ctx.UserId(),
		Remarks: params.Remark,
		ApplyId: int(params.ApplyId),
	}); err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&web.ContactApplyDeclineResponse{})
}

// List 获取联系人申请列表
func (c *Apply) List(ctx *ichat.Context) error {

	list, err := c.ContactApplyService.List(ctx.Ctx(), ctx.UserId())
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

	c.ContactApplyService.ClearApplyUnreadNum(ctx.Ctx(), ctx.UserId())

	return ctx.Success(&web.ContactApplyListResponse{Items: items})
}
