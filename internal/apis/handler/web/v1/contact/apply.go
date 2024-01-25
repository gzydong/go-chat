package contact

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
)

type Apply struct {
	ContactRepo         *repo.Contact
	ContactApplyService service.IContactApplyService
	UserService         service.IUserService
	ContactService      service.IContactService
	MessageService      message.IService
}

// ApplyUnreadNum 获取好友申请未读数
func (c *Apply) ApplyUnreadNum(ctx *core.Context) error {
	return ctx.Success(map[string]any{
		"unread_num": c.ContactApplyService.GetApplyUnreadNum(ctx.Ctx(), ctx.UserId()),
	})
}

// Create 创建联系人申请
func (c *Apply) Create(ctx *core.Context) error {
	in := &web.ContactApplyCreateRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if c.ContactRepo.IsFriend(ctx.Ctx(), uid, int(in.UserId), false) {
		return ctx.Success(nil)
	}

	if err := c.ContactApplyService.Create(ctx.Ctx(), &service.ContactApplyCreateOpt{
		UserId:   ctx.UserId(),
		Remarks:  in.Remark,
		FriendId: int(in.UserId),
	}); err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&web.ContactApplyCreateResponse{})
}

// Accept 同意联系人添加申请
func (c *Apply) Accept(ctx *core.Context) error {
	in := &web.ContactApplyAcceptRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	applyInfo, err := c.ContactApplyService.Accept(ctx.Ctx(), &service.ContactApplyAcceptOpt{
		Remarks: in.Remark,
		ApplyId: int(in.ApplyId),
		UserId:  uid,
	})

	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	_ = c.MessageService.CreatePrivateSysMessage(ctx.Ctx(), message.CreatePrivateSysMessageOption{
		FromId:   uid,
		ToFromId: applyInfo.UserId,
		Content:  "你们已成为好友，可以开始聊天咯！",
	})

	_ = c.MessageService.CreatePrivateSysMessage(ctx.Ctx(), message.CreatePrivateSysMessageOption{
		FromId:   applyInfo.UserId,
		ToFromId: uid,
		Content:  "你们已成为好友，可以开始聊天咯！",
	})

	return ctx.Success(&web.ContactApplyAcceptResponse{})
}

// Decline 拒绝联系人添加申请
func (c *Apply) Decline(ctx *core.Context) error {
	in := &web.ContactApplyDeclineRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ContactApplyService.Decline(ctx.Ctx(), &service.ContactApplyDeclineOpt{
		UserId:  ctx.UserId(),
		Remarks: in.Remark,
		ApplyId: int(in.ApplyId),
	}); err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&web.ContactApplyDeclineResponse{})
}

// List 获取联系人申请列表
func (c *Apply) List(ctx *core.Context) error {

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
