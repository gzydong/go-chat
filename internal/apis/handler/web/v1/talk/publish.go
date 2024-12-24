package talk

import (
	"context"
	"html"

	"github.com/gin-gonic/gin/binding"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/service"
	"go-chat/internal/service/message"
)

var mapping map[string]func(ctx *core.Context) error

type Publish struct {
	AuthService    service.IAuthService
	MessageService message.IService
}

type BaseMessageRequest struct {
	Type     string `json:"type" binding:"required"`            // 消息类型 text:文本消息 image:图片消息 voice:语音消息 video:视频消息 file:文件消息 location:位置消息
	TalkMode int    `json:"talk_mode" binding:"required,gt=0"`  // 对话类型 1:私聊 2:群聊
	ToFromId int    `json:"to_from_id" binding:"required,gt=0"` // 接受者ID (好友ID或者群ID)
	QuoteId  string `json:"quote_id"`                           // 引用的消息ID
}

// Send 发送消息接口
func (c *Publish) Send(ctx *core.Context) error {
	in := &BaseMessageRequest{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:          in.TalkMode,
		UserId:            ctx.UserId(),
		ToFromId:          in.ToFromId,
		IsVerifyGroupMute: true,
	}); err != nil {
		return ctx.Error(err)
	}

	return c.transfer(ctx, in.Type)
}

type onSendTextMessage struct {
	BaseMessageRequest
	Body struct {
		Text     string `json:"text" binding:"required"`
		Mentions []int  `json:"mentions"`
	} `json:"body" binding:"required"`
}

// 文本消息
func (c *Publish) onSendText(ctx *core.Context) error {
	in := &onSendTextMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateTextMessage(ctx.Ctx(), message.CreateTextMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		Content:  html.EscapeString(in.Body.Text),
		QuoteId:  in.QuoteId,
		Mentions: in.Body.Mentions,
	})

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendImageMessage struct {
	BaseMessageRequest
	Body struct {
		Url    string `json:"url" binding:"required"`
		Width  int    `json:"width" binding:"required"`
		Height int    `json:"height" binding:"required"`
		Size   int    `json:"size" binding:"required"`
	} `json:"body" binding:"required"`
}

// 图片消息
func (c *Publish) onSendImage(ctx *core.Context) error {
	in := &onSendImageMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateImageMessage(ctx.Ctx(), message.CreateImageMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		QuoteId:  in.QuoteId,
		Url:      in.Body.Url,
		Width:    in.Body.Width,
		Height:   in.Body.Height,
		Size:     in.Body.Size,
	})

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendVoiceMessage struct {
	BaseMessageRequest
	Body struct {
		Url      string `json:"url" binding:"required"`
		Duration int    `json:"duration" binding:"required"`
		Size     int    `json:"size" binding:"required"`
	} `json:"body" binding:"required"`
}

// 语音消息
func (c *Publish) onSendVoice(ctx *core.Context) error {
	in := &onSendVoiceMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateVoiceMessage(ctx.Ctx(), message.CreateVoiceMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		Url:      in.Body.Url,
		Duration: in.Body.Duration,
		Size:     in.Body.Size,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendVideoMessage struct {
	BaseMessageRequest
	Body struct {
		Url      string `json:"url" binding:"required"`
		Duration int    `json:"duration" binding:"required"`
		Size     int    `json:"size" binding:"required"`
		Cover    string `json:"cover" binding:"required"`
	} `json:"body" binding:"required"`
}

// 视频消息
func (c *Publish) onSendVideo(ctx *core.Context) error {
	in := &onSendVideoMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateVideoMessage(ctx.Ctx(), message.CreateVideoMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		Url:      in.Body.Url,
		Duration: in.Body.Duration,
		Size:     in.Body.Size,
		Cover:    in.Body.Cover,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendFileMessage struct {
	BaseMessageRequest
	Body struct {
		UploadId string `json:"upload_id" binding:"required"`
	} `json:"body" binding:"required"`
}

// 文件消息
func (c *Publish) onSendFile(ctx *core.Context) error {
	in := &onSendFileMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateFileMessage(ctx.Ctx(), message.CreateFileMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		UploadId: in.Body.UploadId,
	})

	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendCodeMessage struct {
	BaseMessageRequest
	Body struct {
		Code string `json:"code" binding:"required"`
		Lang string `json:"lang" binding:"required"`
	} `json:"body" binding:"required"`
}

// 代码消息
func (c *Publish) onSendCode(ctx *core.Context) error {
	in := &onSendCodeMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateCodeMessage(ctx.Ctx(), message.CreateCodeMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		Code:     in.Body.Code,
		Lang:     in.Body.Lang,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendLocationMessage struct {
	BaseMessageRequest
	Body struct {
		Latitude    string `json:"latitude" binding:"required"`
		Longitude   string `json:"longitude" binding:"required"`
		Description string `json:"description" binding:"required"`
	} `json:"body" binding:"required"`
}

// 位置消息
func (c *Publish) onSendLocation(ctx *core.Context) error {
	in := &onSendLocationMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateLocationMessage(ctx.Ctx(), message.CreateLocationMessage{
		TalkMode:    in.TalkMode,
		FromId:      ctx.UserId(),
		ToFromId:    in.ToFromId,
		Longitude:   in.Body.Longitude,
		Latitude:    in.Body.Latitude,
		Description: in.Body.Description,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendForwardMessage struct {
	BaseMessageRequest
	Body struct {
		UserIds  []int    `json:"user_ids"`                   // 好友ID列表
		GroupIds []int    `json:"group_ids"`                  // 群ID列表
		MsgIds   []string `json:"msg_ids" binding:"required"` // 消息ID列表
		Action   int32    `json:"action" binding:"required"`  // 转发模式
	} `json:"body" binding:"required"`
}

// 转发消息
func (c *Publish) onSendForward(ctx *core.Context) error {
	in := &onSendForwardMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	if len(in.Body.MsgIds) == 0 {
		return ctx.InvalidParams("请选择要转发的消息")
	}

	go func() {
		err := c.MessageService.CreateForwardMessage(context.Background(), message.CreateForwardMessage{
			TalkMode: in.TalkMode,
			FromId:   ctx.UserId(),
			ToFromId: in.ToFromId,
			Action:   int(in.Body.Action),
			MsgIds:   in.Body.MsgIds,
			Gids:     in.Body.GroupIds,
			Uids:     in.Body.UserIds,
			UserId:   ctx.UserId(),
		})
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	return ctx.Success(nil)
}

type onSendEmoticonMessage struct {
	BaseMessageRequest
	Body struct {
		EmoticonId int `json:"emoticon_id" binding:"required"`
	} `json:"body" binding:"required"`
}

// 表情消息
func (c *Publish) onSendEmoticon(ctx *core.Context) error {
	in := &onSendEmoticonMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateEmoticonMessage(ctx.Ctx(), message.CreateEmoticonMessage{
		TalkMode:   in.TalkMode,
		FromId:     ctx.UserId(),
		ToFromId:   in.ToFromId,
		EmoticonId: in.Body.EmoticonId,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onSendCardMessage struct {
	BaseMessageRequest
	Body struct {
		UserId int `json:"user_id" binding:"required"`
	} `json:"body" binding:"required"`
}

// 名片消息
func (c *Publish) onSendCard(ctx *core.Context) error {
	in := &onSendCardMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.MessageService.CreateBusinessCardMessage(ctx.Ctx(), message.CreateBusinessCardMessage{
		TalkMode: in.TalkMode,
		FromId:   ctx.UserId(),
		ToFromId: in.ToFromId,
		UserId:   in.Body.UserId,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

type onMixedMessageMessage struct {
	BaseMessageRequest
	Body struct {
		Items []struct {
			Type    int    `json:"type" binding:"required"`
			Content string `json:"content" binding:"required"`
		} `json:"items" binding:"required"`
	} `json:"body" binding:"required"`
}

// 图文消息
func (c *Publish) onMixedMessage(ctx *core.Context) error {
	in := &onMixedMessageMessage{}
	if err := ctx.Context.ShouldBindBodyWith(in, binding.JSON); err != nil {
		return ctx.InvalidParams(err)
	}

	items := make([]message.CreateMixedMessageItem, 0)
	for _, item := range in.Body.Items {
		items = append(items, message.CreateMixedMessageItem{
			Type:    item.Type,
			Content: item.Content,
		})
	}

	err := c.MessageService.CreateMixedMessage(ctx.Ctx(), message.CreateMixedMessage{
		TalkMode:    in.TalkMode,
		FromId:      ctx.UserId(),
		ToFromId:    in.ToFromId,
		QuoteId:     in.QuoteId,
		MessageList: items,
	})
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(nil)
}

func (c *Publish) transfer(ctx *core.Context, typeValue string) error {
	if mapping == nil {
		mapping = make(map[string]func(ctx *core.Context) error)
		mapping["text"] = c.onSendText
		mapping["code"] = c.onSendCode
		mapping["location"] = c.onSendLocation
		mapping["emoticon"] = c.onSendEmoticon
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
