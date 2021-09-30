package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/service"
)

type TalkMessage struct {
	talkMessageService *service.TalkMessageService
}

// Text 发送文本消息
func (c *TalkMessage) Text(ctx *gin.Context) {

}

// Code 发送代码块消息
func (c *TalkMessage) Code(ctx *gin.Context) {

}

// Image 发送图片消息
func (c *TalkMessage) Image(ctx *gin.Context) {

}

// File 发送文件消息
func (c *TalkMessage) File(ctx *gin.Context) {

}

// Vote 发送投票消息
func (c *TalkMessage) Vote(ctx *gin.Context) {

}

// Emoticon 发送表情包消息
func (c *TalkMessage) Emoticon(ctx *gin.Context) {

}

// Forward 发送转发消息
func (c *TalkMessage) Forward(ctx *gin.Context) {

}

// Card 发送用户名片消息
func (c *TalkMessage) Card(ctx *gin.Context) {

}

// Collect 收藏聊天图片
func (c *TalkMessage) Collect(ctx *gin.Context) {

}

// Revoke 撤销聊天记录
func (c *TalkMessage) Revoke(ctx *gin.Context) {

}

// Delete 删除聊天记录
func (c *TalkMessage) Delete(ctx *gin.Context) {

}

// HandleVote 投票处理
func (c *TalkMessage) HandleVote(ctx *gin.Context) {

}
