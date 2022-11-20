package talk

import (
	"github.com/gin-gonic/gin/binding"
	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service"
)

type SendMessage struct {
	auth    *service.TalkAuthService
	message *service.MessageService
}

func NewSendMessage(auth *service.TalkAuthService, message *service.MessageService) *SendMessage {
	return &SendMessage{auth: auth, message: message}
}

type SendBaseMessageRequest struct {
	Type     int       `json:"type" binding:"required,gt=0"`
	Receiver *Receiver `json:"receiver" binding:"required"`
}

// Receiver 接受者信息
type Receiver struct {
	TalkType   int `json:"talk_type" binding:"required,gt=0"`   // 对话类型 1:私聊 2:群聊
	ReceiverId int `json:"receiver_id" binding:"required,gt=0"` // 好友ID或群ID
}

// Send 发送消息接口
func (c *SendMessage) Send(ctx *ichat.Context) error {

	params := &SendBaseMessageRequest{}
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

	return c.transfer(ctx, params.Type)
}

func (c *SendMessage) transfer(ctx *ichat.Context, typeValue int) error {
	switch typeValue {
	case entity.MsgTypeText:
		return c.onSendText(ctx)
	case entity.MsgTypeCode:
		return c.onSendCode(ctx)
	case entity.MsgTypeForward:
		return c.onSendForward(ctx)
	case entity.MsgTypeLocation:
		return c.onSendLocation(ctx)
	case entity.MsgTypeEmoticon:
		return c.onSendEmoticon(ctx)
	case entity.MsgTypeVote:
		return c.onSendVote(ctx)
	default:
		return ctx.InvalidParams("消息类型未定义")
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

// 文件消息
func (c *SendMessage) onSendFile(ctx *ichat.Context) error {

	params := &message.FileMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendFile(ctx.Ctx(), ctx.UserId(), params)
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

// 表情消息
func (c *SendMessage) onSendEmoticon(ctx *ichat.Context) error {

	params := &message.EmoticonMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.message.SendEmoticon(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// 投票消息
func (c *SendMessage) onSendVote(ctx *ichat.Context) error {

	params := &message.VoteMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	if len(params.Options) <= 1 {
		return ctx.InvalidParams("options 选项必须大于1！")
	}

	if len(params.Options) > 6 {
		return ctx.InvalidParams("options 选项不能超过6个！")
	}

	err := c.message.SendVote(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}
