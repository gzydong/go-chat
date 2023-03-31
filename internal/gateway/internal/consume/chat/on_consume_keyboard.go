package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
)

type ConsumeTalkKeyboard struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}

// 键盘输入事件消息
func (h *Handler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {

	var msg ConsumeTalkKeyboard
	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: ", err.Error())
		return
	}

	cids := h.clientStorage.GetUidFromClientIds(ctx, h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(msg.ReceiverID))
	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: "im.message.keyboard",
		Content: map[string]any{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
		},
	})

	socket.Session.Chat.Write(c)
}
