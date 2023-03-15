package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
)

type ConsumeLogin struct {
	Status int `json:"status"`
	UserID int `json:"user_id"`
}

// onConsumeLogin 用户上线或下线消息
func (h *Handler) onConsumeLogin(body []byte) {
	var msg ConsumeLogin

	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeLogin Unmarshal err: ", err.Error())
		return
	}

	ctx := context.Background()
	cids := make([]int64, 0)

	uids := h.contactService.GetContactIds(ctx, msg.UserID)
	sid := h.config.ServerId()
	for _, uid := range uids {
		ids := h.clientStorage.GetUidFromClientIds(ctx, sid, socket.Session.Chat.Name(), strconv.FormatInt(uid, 10))

		cids = append(cids, ids...)
	}

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event:   entity.EventOnlineStatus,
		Content: msg,
	})

	socket.Session.Chat.Write(c)
}
