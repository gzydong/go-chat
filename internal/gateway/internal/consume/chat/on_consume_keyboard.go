package chat

import (
	"context"
	"encoding/json"
	"strconv"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
)

type ConsumeTalkKeyboard struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}

// onConsumeTalkKeyboard 键盘输入事件消息
func (h *Handler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {
	var msg ConsumeTalkKeyboard

	if err := json.Unmarshal(body, &msg); err != nil {
		logger.Error("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: ", err.Error())
		return
	}

	cids := h.clientStorage.GetUidFromClientIds(context.Background(), h.config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(msg.ReceiverID))

	if len(cids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(cids...)
	c.SetMessage(&socket.Message{
		Event: entity.EventTalkKeyboard,
		Content: entity.MapStrAny{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
		},
	})

	socket.Session.Chat.Write(c)
}
