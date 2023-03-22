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

	var data *ConsumeTalkRead
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Error("[ChatSubscribe] onConsumeContactApply Unmarshal err: ", err.Error())
		return
	}

	cids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(data.ReceiverId))

	c := socket.NewSenderContent()
	c.IsAck = true
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: "im.message.read",
		Content: entity.MapStrAny{
			"sender_id":   data.SenderId,
			"receiver_id": data.ReceiverId,
			"ids":         data.Ids,
		},
	})

	socket.Session.Chat.Write(c)
}
