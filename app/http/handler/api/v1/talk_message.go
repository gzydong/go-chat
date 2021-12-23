package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/service"
)

type TalkMessage struct {
	service            *service.TalkMessageService
	talkService        *service.TalkService
	talkRecordsVoteDao *dao.TalkRecordsVoteDao
	forwardService     *service.TalkMessageForwardService
	splitUploadService *service.SplitUploadService
	contactService     *service.ContactService
	groupMemberService *service.GroupMemberService
}

func NewTalkMessageHandler(service *service.TalkMessageService, talkService *service.TalkService, talkRecordsVoteDao *dao.TalkRecordsVoteDao, forwardService *service.TalkMessageForwardService, splitUploadService *service.SplitUploadService, contactService *service.ContactService, groupMemberService *service.GroupMemberService) *TalkMessage {
	return &TalkMessage{service: service, talkService: talkService, talkRecordsVoteDao: talkRecordsVoteDao, forwardService: forwardService, splitUploadService: splitUploadService, contactService: contactService, groupMemberService: groupMemberService}
}

type AuthPermission struct {
	ctx        context.Context
	TalkType   int
	UserId     int
	ReceiverId int
}

// 权限控制
func (c *TalkMessage) permission(prem *AuthPermission) bool {
	// todo 后面需要加缓存
	if prem.TalkType == entity.PrivateChat {
		return c.contactService.Dao().IsFriend(prem.ctx, prem.UserId, prem.ReceiverId, true)
	} else {
		return c.groupMemberService.Dao().IsMember(prem.ReceiverId, prem.UserId, true)
	}
}

// Text 发送文本消息
func (c *TalkMessage) Text(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     auth.GetAuthUserID(ctx),
		ReceiverId: params.ReceiverId,
	}) {
		response.Unauthorized(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendTextMessage(ctx.Request.Context(), uid, params); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Code 发送代码块消息
func (c *TalkMessage) Code(ctx *gin.Context) {
	params := &request.CodeMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendCodeMessage(ctx.Request.Context(), uid, params); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Image 发送图片消息
func (c *TalkMessage) Image(ctx *gin.Context) {
	params := &request.ImageMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		response.InvalidParams(ctx, "image 字段必传！")
		return
	}

	if !slice.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif"}) {
		response.InvalidParams(ctx, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(ctx, "上传文件大小不能超过5M！")
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendImageMessage(ctx.Request.Context(), uid, params, file); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// File 发送文件消息
func (c *TalkMessage) File(ctx *gin.Context) {
	params := &request.FileMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	file, err := c.splitUploadService.Dao().GetFile(uid, params.UploadId)
	if err != nil {
		response.BusinessError(ctx, "文件信息不存在！")
		return
	}

	if err := c.service.SendFileMessage(ctx.Request.Context(), uid, params, file); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Vote 发送投票消息
func (c *TalkMessage) Vote(ctx *gin.Context) {
	params := &request.VoteMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if len(params.Options) <= 1 {
		response.InvalidParams(ctx, "options 选项必须大于1！")
		return
	}

	if len(params.Options) > 6 {
		response.InvalidParams(ctx, "options 选项不能超过6个！")
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   entity.GroupChat,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendVoteMessage(ctx.Request.Context(), uid, params); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Emoticon 发送表情包消息
func (c *TalkMessage) Emoticon(ctx *gin.Context) {
	params := &request.EmoticonMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendEmoticonMessage(ctx.Request.Context(), uid, params); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Forward 发送转发消息
func (c *TalkMessage) Forward(ctx *gin.Context) {
	params := &request.ForwardMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if params.ReceiveGroupIds == "" && params.ReceiveUserIds == "" {
		response.InvalidParams(ctx, "receive_user_ids 和 receive_group_ids 不能都为空！")
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	forward := &service.ForwardParams{
		Mode:       params.ForwardMode,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
		TalkType:   params.TalkType,
		RecordsIds: slice.ParseIds(params.RecordsIds),
		UserIds:    slice.ParseIds(params.ReceiveUserIds),
		GroupIds:   slice.ParseIds(params.ReceiveGroupIds),
	}

	if err := c.forwardService.SendForwardMessage(ctx.Request.Context(), forward); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Card 发送用户名片消息
func (c *TalkMessage) Card(ctx *gin.Context) {
	params := &request.CardMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	// c.service.SendCardMessage(ctx.Request.Context(), params)

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Collect 收藏聊天图片
func (c *TalkMessage) Collect(ctx *gin.Context) {
	params := &request.CollectMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.talkService.CollectRecord(ctx.Request.Context(), auth.GetAuthUserID(ctx), params.RecordId)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{}, "已收藏！")
}

// Revoke 撤销聊天记录
func (c *TalkMessage) Revoke(ctx *gin.Context) {
	params := &request.RevokeMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.SendRevokeRecordMessage(ctx.Request.Context(), auth.GetAuthUserID(ctx), params.RecordId)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// Delete 删除聊天记录
func (c *TalkMessage) Delete(ctx *gin.Context) {
	params := &request.DeleteMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.talkService.RemoveRecords(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}

// HandleVote 投票处理
func (c *TalkMessage) HandleVote(ctx *gin.Context) {
	params := &request.VoteMessageHandleRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	vid, err := c.service.VoteHandle(ctx.Request.Context(), auth.GetAuthUserID(ctx), params)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	// 获取同级数据
	res, _ := c.talkRecordsVoteDao.GetVoteStatistics(ctx.Request.Context(), vid)

	response.Success(ctx, res, "消息推送成功！")
}

// Location 发送位置消息
func (c *TalkMessage) Location(ctx *gin.Context) {
	params := &request.LocationMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := auth.GetAuthUserID(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendLocationMessage(ctx.Request.Context(), uid, params); err != nil {
		response.Success(ctx, gin.H{}, "消息推送失败！")
		return
	}

	response.Success(ctx, gin.H{}, "消息推送成功！")
}
