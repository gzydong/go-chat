package talk

import (
	"github.com/gin-gonic/gin/binding"
	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service"
)

type SendMessage struct {
	auth    *service.TalkAuthService
	message *service.MessageService
}

func (c *SendMessage) Send(ctx *ichat.Context) error {

	params := &web.SendBaseMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	// 权限验证
	if err := c.auth.IsAuth(ctx.Ctx(), &service.TalkAuthOption{
		TalkType:   params.Receiver.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.Receiver.ReceiverId,
	}); err != nil {
		return ctx.BusinessError(err.Error())
	}

	switch params.Type {
	case entity.MsgTypeText:
		return c.onSendText(ctx)
	case entity.MsgTypeCode:
		return c.onSendCode(ctx)
	case entity.MsgTypeForward:
		return c.onSendForward(ctx)
	case entity.MsgTypeLocation:
		return c.onSendLocation(ctx)
	default:
		return ctx.InvalidParams("消息类型不能为空")
	}
}

// 文本消息
func (c *SendMessage) onSendText(ctx *ichat.Context) error {

	params := &message.TextMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendText(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// 图片消息
func (c *SendMessage) onSendImage(ctx *ichat.Context) error {

	params := &message.ImageMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendImage(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// 代码消息
func (c *SendMessage) onSendCode(ctx *ichat.Context) error {

	params := &message.CodeMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendCode(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// 位置消息
func (c *SendMessage) onSendLocation(ctx *ichat.Context) error {

	params := &message.LocationMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendLocation(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// 转发消息
func (c *SendMessage) onSendForward(ctx *ichat.Context) error {

	params := &message.ForwardMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendForward(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}
