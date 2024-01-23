package talk

import (
	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service"
)

type Message struct {
	TalkService    service.ITalkService
	AuthService    service.IAuthService
	MessageService service.IMessageService
	Filesystem     filesystem.IFilesystem
}

type CollectMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required"`
	MsgId      string `form:"msg_id" json:"msg_id" binding:"required"`
}

// Collect 收藏聊天图片
func (c *Message) Collect(ctx *ichat.Context) error {
	params := &CollectMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkService.Collect(ctx.Ctx(), &service.CollectOpt{
		UserId:     ctx.UserId(),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		MsgId:      params.MsgId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type RevokeMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required"`
	MsgId      string `form:"msg_id" json:"msg_id" binding:"required"`
}

// Revoke 撤销聊天记录
func (c *Message) Revoke(ctx *ichat.Context) error {
	params := &RevokeMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.MessageService.Revoke(ctx.Ctx(), ctx.UserId(), params.MsgId); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type DeleteMessageRequest struct {
	TalkType   int      `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2"`
	ReceiverId int      `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0"`
	MsgIds     []string `form:"msg_ids" json:"msg_ids" binding:"required"`
}

// Delete 删除聊天记录
func (c *Message) Delete(ctx *ichat.Context) error {
	params := &DeleteMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkService.DeleteRecordList(ctx.Ctx(), &service.RemoveRecordListOpt{
		UserId:     ctx.UserId(),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		MsgIds:     params.MsgIds,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type VoteMessageRequest struct {
	ReceiverId int      `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0"`
	Mode       int      `form:"mode" json:"mode" binding:"oneof=0 1"`
	Anonymous  int      `form:"anonymous" json:"anonymous" binding:"oneof=0 1"`
	Title      string   `form:"title" json:"title" binding:"required"`
	Options    []string `form:"options" json:"options"`
}

// Vote 发送投票消息
func (c *Message) Vote(ctx *ichat.Context) error {
	params := &VoteMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if len(params.Options) <= 1 {
		return ctx.InvalidParams("options 选项必须大于1！")
	}

	if len(params.Options) > 6 {
		return ctx.InvalidParams("options 选项不能超过6个！")
	}

	uid := ctx.UserId()
	if err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:          entity.ChatGroupMode,
		UserId:            uid,
		ReceiverId:        params.ReceiverId,
		IsVerifyGroupMute: true,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.MessageService.SendVote(ctx.Ctx(), uid, &message.VoteMessageRequest{
		Mode:      int32(params.Mode),
		Title:     params.Title,
		Options:   params.Options,
		Anonymous: int32(params.Anonymous),
		Receiver: &message.MessageReceiver{
			TalkType:   entity.ChatGroupMode,
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type VoteMessageHandleRequest struct {
	VoteId  int    `form:"vote_id" json:"vote_id" binding:"required,numeric,gt=0"`
	Options string `form:"options" json:"options" binding:"required"`
}

// SubmitVote 提交投票
func (c *Message) SubmitVote(ctx *ichat.Context) error {
	params := &VoteMessageHandleRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	data, err := c.MessageService.Vote(ctx.Ctx(), ctx.UserId(), params.VoteId, params.Options)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(data)
}
