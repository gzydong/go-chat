package talk

import (
	"bytes"
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/service"
)

type Message struct {
	talkService    *service.TalkService
	authService    *service.AuthService
	messageService *service.MessageService
	filesystem     *filesystem.Filesystem
}

func NewMessage(talkService *service.TalkService, talkAuthService *service.AuthService, messageService *service.MessageService, filesystem *filesystem.Filesystem) *Message {
	return &Message{talkService: talkService, authService: talkAuthService, messageService: messageService, filesystem: filesystem}
}

type AuthorityOption struct {
	TalkType   int // 对话类型
	UserId     int // 发送者ID
	ReceiverId int // 接收者ID
}

type TextMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Text       string `form:"text" json:"text" binding:"required,max=3000" label:"text"`
}

// Text 发送文本消息
func (c *Message) Text(ctx *ichat.Context) error {

	params := &TextMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendText(ctx.Ctx(), uid, &message.TextMessageRequest{
		Content: params.Text,
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type CodeMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Lang       string `form:"lang" json:"lang" binding:"required"`
	Code       string `form:"code" json:"code" binding:"required,max=65535"`
}

// Code 发送代码块消息
func (c *Message) Code(ctx *ichat.Context) error {

	params := &CodeMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendCode(ctx.Ctx(), uid, &message.CodeMessageRequest{
		Lang: params.Lang,
		Code: params.Code,
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type ImageMessageRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
}

// Image 发送图片消息
func (c *Message) Image(ctx *ichat.Context) error {

	params := &ImageMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	file, err := ctx.Context.FormFile("image")
	if err != nil {
		return ctx.InvalidParams("image 字段必传！")
	}

	if !sliceutil.Include(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif", "webp"}) {
		return ctx.InvalidParams("上传文件格式不正确,仅支持 png、jpg、jpeg、gif 及 webp")
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		return ctx.InvalidParams("上传文件大小不能超过5M！")
	}

	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return err
	}

	ext := strutil.FileSuffix(file.Filename)

	meta := utils.ReadImageMeta(bytes.NewReader(stream))

	filePath := fmt.Sprintf("public/media/image/talk/%s/%s", timeutil.DateNumber(), strutil.GenImageName(ext, meta.Width, meta.Height))

	if err := c.filesystem.Default.Write(stream, filePath); err != nil {
		return err
	}

	if err := c.messageService.SendImage(ctx.Ctx(), ctx.UserId(), &message.ImageMessageRequest{
		Url:    c.filesystem.Default.PublicUrl(filePath),
		Width:  int32(meta.Width),
		Height: int32(meta.Height),
		Size:   int32(file.Size),
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type FileMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	UploadId   string `form:"upload_id" json:"upload_id" binding:"required"`
}

// File 发送文件消息
func (c *Message) File(ctx *ichat.Context) error {

	params := &FileMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendFile(ctx.Ctx(), uid, &message.FileMessageRequest{
		UploadId: params.UploadId,
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type VoteMessageRequest struct {
	ReceiverId int      `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
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
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   entity.ChatGroupMode,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendVote(ctx.Ctx(), uid, &message.VoteMessageRequest{
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

type EmoticonMessageRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	EmoticonId int `form:"emoticon_id" json:"emoticon_id" binding:"required,numeric,gt=0"`
}

// Emoticon 发送表情包消息
func (c *Message) Emoticon(ctx *ichat.Context) error {

	params := &EmoticonMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendEmoticon(ctx.Ctx(), uid, &message.EmoticonMessageRequest{
		EmoticonId: int32(params.EmoticonId),
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type ForwardMessageRequest struct {
	TalkType        int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId      int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	ForwardMode     int    `form:"forward_mode" json:"forward_mode" binding:"required,oneof=1 2"`
	RecordsIds      string `form:"records_ids" json:"records_ids" binding:"required,ids"`
	ReceiveUserIds  string `form:"receive_user_ids" json:"receive_user_ids" binding:"ids"`
	ReceiveGroupIds string `form:"receive_group_ids" json:"receive_group_ids" binding:"ids"`
}

// Forward 发送转发消息
func (c *Message) Forward(ctx *ichat.Context) error {

	params := &ForwardMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.ReceiveGroupIds == "" && params.ReceiveUserIds == "" {
		return ctx.InvalidParams("receive_user_ids 和 receive_group_ids 不能都为空！")
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	data := &message.ForwardMessageRequest{
		Mode:       int32(params.ForwardMode),
		MessageIds: make([]int32, 0),
		Gids:       make([]int32, 0),
		Uids:       make([]int32, 0),
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}

	for _, id := range sliceutil.ParseIds(params.RecordsIds) {
		data.MessageIds = append(data.MessageIds, int32(id))
	}

	for _, id := range sliceutil.ParseIds(params.ReceiveUserIds) {
		data.Uids = append(data.Uids, int32(id))
	}

	for _, id := range sliceutil.ParseIds(params.ReceiveGroupIds) {
		data.Gids = append(data.Gids, int32(id))
	}

	if err := c.messageService.SendForward(ctx.Ctx(), uid, data); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type CardMessageRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
}

// Card 发送用户名片消息
func (c *Message) Card(ctx *ichat.Context) error {

	params := &CardMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendBusinessCard(ctx.Ctx(), uid, &message.CardMessageRequest{
		UserId: int32(params.ReceiverId),
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type CollectMessageRequest struct {
	RecordId int `form:"record_id" json:"record_id" binding:"required,numeric,gt=0" label:"record_id"`
}

// Collect 收藏聊天图片
func (c *Message) Collect(ctx *ichat.Context) error {

	params := &CollectMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.talkService.Collect(ctx.Ctx(), ctx.UserId(), params.RecordId); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type RevokeMessageRequest struct {
	RecordId int `form:"record_id" json:"record_id" binding:"required,numeric,gt=0" label:"record_id"`
}

// Revoke 撤销聊天记录
func (c *Message) Revoke(ctx *ichat.Context) error {

	params := &RevokeMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.messageService.Revoke(ctx.Ctx(), ctx.UserId(), params.RecordId); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type DeleteMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	RecordIds  string `form:"record_id" json:"record_id" binding:"required,ids" label:"record_id"`
}

// Delete 删除聊天记录
func (c *Message) Delete(ctx *ichat.Context) error {

	params := &DeleteMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.talkService.DeleteRecordList(ctx.Ctx(), &service.RemoveRecordListOpt{
		UserId:     ctx.UserId(),
		TalkType:   params.TalkType,
		ReceiverId: params.ReceiverId,
		RecordIds:  params.RecordIds,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}

type VoteMessageHandleRequest struct {
	RecordId int    `form:"record_id" json:"record_id" binding:"required,gt=0"`
	Options  string `form:"options" json:"options" binding:"required"`
}

// HandleVote 投票处理
func (c *Message) HandleVote(ctx *ichat.Context) error {

	params := &VoteMessageHandleRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	data, err := c.messageService.Vote(ctx.Ctx(), ctx.UserId(), params.RecordId, params.Options)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(data)
}

type LocationMessageRequest struct {
	TalkType   int    `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2" label:"talk_type"`
	ReceiverId int    `form:"receiver_id" json:"receiver_id" binding:"required,numeric,gt=0" label:"receiver_id"`
	Longitude  string `form:"longitude" json:"longitude" binding:"required,numeric"`
	Latitude   string `form:"latitude" json:"latitude" binding:"required,numeric"`
}

// Location 发送位置消息
func (c *Message) Location(ctx *ichat.Context) error {

	params := &LocationMessageRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if err := c.messageService.SendLocation(ctx.Ctx(), uid, &message.LocationMessageRequest{
		Longitude:   params.Longitude,
		Latitude:    params.Latitude,
		Description: "", // todo 需完善
		Receiver: &message.MessageReceiver{
			TalkType:   int32(params.TalkType),
			ReceiverId: int32(params.ReceiverId),
		},
	}); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(nil)
}
