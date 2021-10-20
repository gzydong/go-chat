package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/service"
)

type TalkMessage struct {
	TalkMessageService *service.TalkMessageService
}

func NewTalkMessageHandler(
	talkMessage *service.TalkMessageService,
) *TalkMessage {
	return &TalkMessage{
		TalkMessageService: talkMessage,
	}
}

// Text 发送文本消息
func (c *TalkMessage) Text(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendTextMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Code 发送代码块消息
func (c *TalkMessage) Code(ctx *gin.Context) {
	params := &request.CodeMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendCodeMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Image 发送图片消息
func (c *TalkMessage) Image(ctx *gin.Context) {
	params := &request.ImageMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendImageMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// File 发送文件消息
func (c *TalkMessage) File(ctx *gin.Context) {
	params := &request.FileMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendFileMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Vote 发送投票消息
func (c *TalkMessage) Vote(ctx *gin.Context) {
	params := &request.VoteMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendVoteMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Emoticon 发送表情包消息
func (c *TalkMessage) Emoticon(ctx *gin.Context) {
	params := &request.EmoticonMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendEmoticonMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Forward 发送转发消息
func (c *TalkMessage) Forward(ctx *gin.Context) {
	params := &request.ForwardMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendForwardMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Card 发送用户名片消息
func (c *TalkMessage) Card(ctx *gin.Context) {
	params := &request.CardMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendCardMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Collect 收藏聊天图片
func (c *TalkMessage) Collect(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendTextMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Revoke 撤销聊天记录
func (c *TalkMessage) Revoke(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendTextMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Delete 删除聊天记录
func (c *TalkMessage) Delete(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendTextMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// HandleVote 投票处理
func (c *TalkMessage) HandleVote(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	c.TalkMessageService.SendTextMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}
