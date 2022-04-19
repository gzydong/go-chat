package talk

import (
	"context"

	"github.com/gin-gonic/gin"

	"go-chat/internal/dao"
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/service"
)

type Message struct {
	service            *service.TalkMessageService
	talkService        *service.TalkService
	talkRecordsVoteDao *dao.TalkRecordsVoteDao
	forwardService     *service.TalkMessageForwardService
	splitUploadService *service.SplitUploadService
	contactService     *service.ContactService
	groupMemberService *service.GroupMemberService
}

func NewTalkMessageHandler(service *service.TalkMessageService, talkService *service.TalkService, talkRecordsVoteDao *dao.TalkRecordsVoteDao, forwardService *service.TalkMessageForwardService, splitUploadService *service.SplitUploadService, contactService *service.ContactService, groupMemberService *service.GroupMemberService) *Message {
	return &Message{service: service, talkService: talkService, talkRecordsVoteDao: talkRecordsVoteDao, forwardService: forwardService, splitUploadService: splitUploadService, contactService: contactService, groupMemberService: groupMemberService}
}

type AuthPermission struct {
	ctx        context.Context
	TalkType   int
	UserId     int
	ReceiverId int
}

// 权限控制
func (c *Message) permission(prem *AuthPermission) bool {
	if prem.TalkType == entity.ChatPrivateMode {
		return c.contactService.Dao().IsFriend(prem.ctx, prem.UserId, prem.ReceiverId, true)
	} else {
		return c.groupMemberService.Dao().IsMember(prem.ReceiverId, prem.UserId, true)
	}
}

// Text 发送文本消息
func (c *Message) Text(ctx *gin.Context) {
	params := &request.TextMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     jwtutil.GetUid(ctx),
		ReceiverId: params.ReceiverId,
	}) {
		response.Unauthorized(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendTextMessage(ctx.Request.Context(), &service.TextMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		Text:       params.Text,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Code 发送代码块消息
func (c *Message) Code(ctx *gin.Context) {
	params := &request.CodeMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendCodeMessage(ctx.Request.Context(), &service.CodeMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		Lang:       params.Lang,
		Code:       params.Code,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Image 发送图片消息
func (c *Message) Image(ctx *gin.Context) {
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

	if !sliceutil.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif"}) {
		response.InvalidParams(ctx, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(ctx, "上传文件大小不能超过5M！")
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendImageMessage(ctx.Request.Context(), &service.ImageMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		File:       file,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// File 发送文件消息
func (c *Message) File(ctx *gin.Context) {
	params := &request.FileMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendFileMessage(ctx.Request.Context(), &service.FileMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		UploadId:   params.UploadId,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Vote 发送投票消息
func (c *Message) Vote(ctx *gin.Context) {
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

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   entity.ChatGroupMode,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendVoteMessage(ctx.Request.Context(), &service.VoteMessageOpts{
		UserId:     uid,
		ReceiverId: params.ReceiverId,
		Mode:       params.Mode,
		Title:      params.Title,
		Options:    params.Options,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Emoticon 发送表情包消息
func (c *Message) Emoticon(ctx *gin.Context) {
	params := &request.EmoticonMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendEmoticonMessage(ctx.Request.Context(), &service.EmoticonMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		EmoticonId: params.EmoticonId,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Forward 发送转发消息
func (c *Message) Forward(ctx *gin.Context) {
	params := &request.ForwardMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if params.ReceiveGroupIds == "" && params.ReceiveUserIds == "" {
		response.InvalidParams(ctx, "receive_user_ids 和 receive_group_ids 不能都为空！")
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	forward := &service.TalkForwardOpts{
		Mode:       params.ForwardMode,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
		TalkType:   params.TalkType,
		RecordsIds: sliceutil.ParseIds(params.RecordsIds),
		UserIds:    sliceutil.ParseIds(params.ReceiveUserIds),
		GroupIds:   sliceutil.ParseIds(params.ReceiveGroupIds),
	}

	if err := c.forwardService.SendForwardMessage(ctx.Request.Context(), forward); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Card 发送用户名片消息
func (c *Message) Card(ctx *gin.Context) {
	params := &request.CardMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	// todo SendCardMessage
	if err := c.service.SendCardMessage(ctx.Request.Context(), &service.CardMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		ContactId:  0,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Collect 收藏聊天图片
func (c *Message) Collect(ctx *gin.Context) {
	params := &request.CollectMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkService.CollectRecord(ctx.Request.Context(), jwtutil.GetUid(ctx), params.RecordId); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Revoke 撤销聊天记录
func (c *Message) Revoke(ctx *gin.Context) {
	params := &request.RevokeMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.SendRevokeRecordMessage(ctx.Request.Context(), jwtutil.GetUid(ctx), params.RecordId); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// Delete 删除聊天记录
func (c *Message) Delete(ctx *gin.Context) {
	params := &request.DeleteMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.talkService.RemoveRecords(ctx.Request.Context(), &service.TalkMessageDeleteOpts{
		UserId:     jwtutil.GetUid(ctx),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		RecordIds:  params.RecordIds,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}

// HandleVote 投票处理
func (c *Message) HandleVote(ctx *gin.Context) {
	params := &request.VoteMessageHandleRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	vid, err := c.service.VoteHandle(ctx.Request.Context(), &service.VoteMessageHandleOpts{
		UserId:   jwtutil.GetUid(ctx),
		RecordId: params.RecordId,
		Options:  params.Options,
	})
	if err != nil {
		response.BusinessError(ctx, err)
	} else {
		res, _ := c.talkRecordsVoteDao.GetVoteStatistics(ctx.Request.Context(), vid)

		response.Success(ctx, res)
	}
}

// Location 发送位置消息
func (c *Message) Location(ctx *gin.Context) {
	params := &request.LocationMessageRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	uid := jwtutil.GetUid(ctx)

	if !c.permission(&AuthPermission{
		ctx:        ctx.Request.Context(),
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}) {
		response.BusinessError(ctx, "无权限访问！")
		return
	}

	if err := c.service.SendLocationMessage(ctx.Request.Context(), &service.LocationMessageOpts{
		UserId:     uid,
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		Longitude:  params.Longitude,
		Latitude:   params.Latitude,
	}); err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil)
	}
}
