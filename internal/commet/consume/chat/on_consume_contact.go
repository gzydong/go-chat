package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
)

type ConsumeContactStatus struct {
	Status int `json:"status"`
	UserId int `json:"user_id"`
}

// 用户上线或下线消息
func (h *Handler) onConsumeContactStatus(ctx context.Context, body []byte) {

	var in ConsumeContactStatus
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactStatus Unmarshal err: %s", err.Error())
		return
	}

	contactIds := h.ContactService.GetContactIds(ctx, in.UserId)

	if isOk, _ := h.OrganizeRepo.IsQiyeMember(ctx, in.UserId); isOk {
		ids, _ := h.OrganizeRepo.GetMemberIds(ctx)
		contactIds = append(contactIds, ids...)
	}

	clientIds := make([]int64, 0)

	for _, uid := range sliceutil.Unique(contactIds) {
		ids := h.ClientStorage.GetUidFromClientIds(ctx, h.Config.ServerId(), socket.Session.Chat.Name(), strconv.FormatInt(uid, 10))
		if len(ids) > 0 {
			clientIds = append(clientIds, ids...)
		}
	}

	if len(clientIds) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventContactStatus, in)

	socket.Session.Chat.Write(c)
}

type ConsumeContactApply struct {
	ApplyId int `json:"apply_id"`
	Type    int `json:"type"`
}

// 好友申请消息
func (h *Handler) onConsumeContactApply(ctx context.Context, body []byte) {

	var in ConsumeContactApply
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactApply Unmarshal err: %s", err.Error())
		return
	}

	var apply model.ContactApply
	if err := h.Source.Db().First(&apply, in.ApplyId).Error; err != nil {
		return
	}

	clientIds := h.ClientStorage.GetUidFromClientIds(ctx, h.Config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(apply.FriendId))
	if len(clientIds) == 0 {
		return
	}

	var user model.Users
	if err := h.Source.Db().First(&user, apply.FriendId).Error; err != nil {
		return
	}

	data := map[string]any{}
	data["sender_id"] = apply.UserId
	data["receiver_id"] = apply.FriendId
	data["remark"] = apply.Remark
	data["friend"] = map[string]any{
		"nickname":   user.Nickname,
		"remark":     apply.Remark,
		"created_at": timeutil.FormatDatetime(apply.CreatedAt),
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventContactApply, data)

	socket.Session.Chat.Write(c)
}
