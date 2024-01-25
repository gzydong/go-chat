package talk

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/service"
)

type Message struct {
	TalkService service.ITalkService
	AuthService service.IAuthService
	Filesystem  filesystem.IFilesystem
}

type RevokeMessageRequest struct {
	TalkMode int    `form:"talk_mode" json:"talk_mode" binding:"required"`
	ToFromId int    `form:"to_from_id" json:"to_from_id" binding:"required"`
	MsgId    string `form:"msg_id" json:"msg_id" binding:"required"`
}

// Revoke 撤销聊天记录
func (c *Message) Revoke(ctx *core.Context) error {
	in := &RevokeMessageRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkService.Revoke(ctx.Ctx(), &service.TalkRevokeOption{
		UserId:   ctx.UserId(),
		TalkMode: in.TalkMode,
		MsgId:    in.MsgId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(map[string]any{})
}

type DeleteMessageRequest struct {
	TalkMode int      `form:"talk_mode" json:"talk_mode" binding:"required"`
	ToFromId int      `form:"to_from_id" json:"to_from_id" binding:"required"`
	MsgIds   []string `form:"msg_ids" json:"msg_ids" binding:"required"`
}

// Delete 删除聊天记录
func (c *Message) Delete(ctx *core.Context) error {
	in := &DeleteMessageRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.TalkService.DeleteRecord(ctx.Ctx(), &service.TalkDeleteRecordOption{
		UserId:   ctx.UserId(),
		TalkMode: in.TalkMode,
		ToFromId: in.ToFromId,
		MsgIds:   in.MsgIds,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}
