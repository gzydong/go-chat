package talk

import (
	"github.com/gin-gonic/gin/binding"
	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service"
)

var mapping map[string]func(ctx *ichat.Context) error

type Publish struct {
	AuthService    service.IAuthService
	MessageService service.IMessageService
}

type PublishBaseMessageRequest struct {
	Type     string `json:"type" binding:"required"`
	Receiver struct {
		TalkType   int `json:"talk_type" binding:"required,gt=0"`   // 对话类型 1:私聊 2:群聊
		ReceiverId int `json:"receiver_id" binding:"required,gt=0"` // 好友ID或群ID
	} `json:"receiver" binding:"required"`
}

// Publish 发送消息接口
func (c *Publish) Publish(ctx *ichat.Context) error {
	params := &PublishBaseMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:          params.Receiver.TalkType,
		UserId:            ctx.UserId(),
		ReceiverId:        params.Receiver.ReceiverId,
		IsVerifyGroupMute: true,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return c.transfer(ctx, params.Type)
}

// 文本消息
func (c *Publish) onSendText(ctx *ichat.Context) error {

	params := &message.TextMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendText(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 图片消息
func (c *Publish) onSendImage(ctx *ichat.Context) error {

	params := &message.ImageMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendImage(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 语音消息
func (c *Publish) onSendVoice(ctx *ichat.Context) error {

	params := &message.VoiceMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendVoice(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 视频消息
func (c *Publish) onSendVideo(ctx *ichat.Context) error {

	params := &message.VideoMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendVideo(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 文件消息
func (c *Publish) onSendFile(ctx *ichat.Context) error {

	params := &message.FileMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendFile(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 代码消息
func (c *Publish) onSendCode(ctx *ichat.Context) error {

	params := &message.CodeMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendCode(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 位置消息
func (c *Publish) onSendLocation(ctx *ichat.Context) error {

	params := &message.LocationMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendLocation(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 转发消息
func (c *Publish) onSendForward(ctx *ichat.Context) error {

	params := &message.ForwardMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendForward(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 表情消息
func (c *Publish) onSendEmoticon(ctx *ichat.Context) error {

	params := &message.EmoticonMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendEmoticon(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 投票消息
func (c *Publish) onSendVote(ctx *ichat.Context) error {

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

	err := c.MessageService.SendVote(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 名片消息
func (c *Publish) onSendCard(ctx *ichat.Context) error {

	params := &message.CardMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendBusinessCard(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

// 图文消息
func (c *Publish) onMixedMessage(ctx *ichat.Context) error {

	params := &message.MixedMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(params, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.SendMixedMessage(ctx.Ctx(), ctx.UserId(), params)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

func (c *Publish) transfer(ctx *ichat.Context, typeValue string) error {
	if mapping == nil {
		mapping = make(map[string]func(ctx *ichat.Context) error)
		mapping["text"] = c.onSendText
		mapping["code"] = c.onSendCode
		mapping["location"] = c.onSendLocation
		mapping["emoticon"] = c.onSendEmoticon
		mapping["vote"] = c.onSendVote
		mapping["image"] = c.onSendImage
		mapping["voice"] = c.onSendVoice
		mapping["video"] = c.onSendVideo
		mapping["file"] = c.onSendFile
		mapping["card"] = c.onSendCard
		mapping["forward"] = c.onSendForward
		mapping["mixed"] = c.onMixedMessage
	}

	if call, ok := mapping[typeValue]; ok {
		return call(ctx)
	}

	return nil
}
