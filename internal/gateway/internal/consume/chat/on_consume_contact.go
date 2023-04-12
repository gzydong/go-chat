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
	UserID int `json:"user_id"`
}

// 用户上线或下线消息
func (h *Handler) onConsumeContactStatus(ctx context.Context, body []byte) {

	var data ConsumeContactStatus
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("[ChatSubscribe] onConsumeContactStatus Unmarshal err: ", err.Error())
		return
	}

	cids := make([]int64, 0)

	uids := h.contactService.GetContactIds(ctx, data.UserID)

	isQiye, _ := h.organize.IsQiyeMember(ctx, data.UserID)
	if isQiye {
		mids, _ := h.organize.GetMemberIds(ctx)
		uids = append(uids, mids...)
	}

	for _, uid := range sliceutil.Unique(uids) {
		ids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.FormatInt(uid, 10))
		if len(ids) > 0 {
			cids = append(cids, ids...)
		}
	}

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(entity.PushEventContactStatus, data)

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
		logger.Error("[ChatSubscribe] onConsumeContactApply Unmarshal err: ", err.Error())
		return
	}

	var apply *model.ContactApply
	if err := h.contactService.Db().First(&apply, in.ApplyId).Error; err != nil {
		return
	}

	clientIds := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(apply.FriendId))
	if len(clientIds) == 0 {
		return
	}

	var user *model.Users
	if err := h.contactService.Db().First(&user, apply.FriendId).Error; err != nil {
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
