package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
)

type ConsumeTalkRead struct {
	SenderId   int   `json:"sender_id"`
	ReceiverId int   `json:"receiver_id"`
	Ids        []int `json:"ids"`
}

// 消息已读事件
func (h *Handler) onConsumeTalkRead(ctx context.Context, body []byte) {

	var in ConsumeTalkRead
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Error("[ChatSubscribe] onConsumeContactApply Unmarshal err: ", err.Error())
		return
	}

	clientIds := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(in.ReceiverId))
	if len(clientIds) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetAck(true)
	c.SetReceive(clientIds...)
	c.SetMessage(entity.PushEventImMessageRead, map[string]any{
		"sender_id":   in.SenderId,
		"receiver_id": in.ReceiverId,
		"ids":         in.Ids,
	})

	socket.Session.Chat.Write(c)
}
