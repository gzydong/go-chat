package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
)

type ConsumeContactApply struct {
	ApplyId int `json:"apply_id"`
	Type    int `json:"type"`
}

// nolint onConsumeContactApply 好友申请消息
func (h *Handler) onConsumeContactApply(ctx context.Context, body []byte) {

	var msg ConsumeContactApply
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeContactApply Unmarshal err: ", err.Error())
		return
	}

	apply := &model.ContactApply{}
	if err := h.contactService.Db().First(&apply, msg.ApplyId).Error; err != nil {
		return
	}

	cids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(apply.FriendId))
	if len(cids) == 0 {
		return
	}

	user := &model.Users{}
	if err := h.contactService.Db().First(&user, apply.FriendId).Error; err != nil {
		return
	}

	data := entity.MapStrAny{}
	data["sender_id"] = apply.UserId
	data["receiver_id"] = apply.FriendId
	data["remark"] = apply.Remark
	data["friend"] = entity.MapStrAny{
		"nickname":   user.Nickname,
		"remark":     apply.Remark,
		"created_at": timeutil.FormatDatetime(apply.CreatedAt),
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event:   entity.EventContactApply,
		Content: data,
	})

	socket.Session.Chat.Write(c)
}
