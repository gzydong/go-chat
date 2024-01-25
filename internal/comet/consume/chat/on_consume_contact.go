package chat

import (
	"context"
	"encoding/json"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/server"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
	"time"
)

// 用户上线或下线消息
func (h *Handler) onConsumeContactStatus(ctx context.Context, body []byte) {
	var in entity.SubEventContactStatusPayload
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
	sid := server.ID()
	for _, uid := range sliceutil.Unique(contactIds) {
		ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), int(uid))
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

// 好友申请消息
func (h *Handler) onConsumeContactApply(ctx context.Context, body []byte) {
	var in entity.SubEventContactApplyPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeContactApply Unmarshal err: %s", err.Error())
		return
	}

	var apply model.ContactApply
	if err := h.Source.Db().First(&apply, in.ApplyId).Error; err != nil {
		return
	}

	clientIds, _ := h.ClientConnectService.GetUidFromClientIds(ctx, server.ID(), socket.Session.Chat.Name(), apply.FriendId)
	if len(clientIds) == 0 {
		return
	}

	var user model.Users
	if err := h.Source.Db().First(&user, apply.UserId).Error; err != nil {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventContactApply, entity.ImContactApplyPayload{
		UserId:    user.Id,
		Nickname:  user.Nickname,
		Remark:    apply.Remark,
		ApplyTime: apply.CreatedAt.Format(time.DateTime),
	})

	socket.Session.Chat.Write(c)
}
