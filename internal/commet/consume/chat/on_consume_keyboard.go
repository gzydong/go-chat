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

// 键盘输入事件消息
func (h *Handler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {

	var in ConsumeTalkKeyboard
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: %s", err.Error())
		return
	}

	ids := h.ClientStorage.GetUidFromClientIds(ctx, h.Config.ServerId(), socket.Session.Chat.Name(), strconv.Itoa(in.ReceiverID))
	if len(ids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(ids...)
	c.SetMessage(entity.PushEventImMessageKeyboard, map[string]any{
		"sender_id":   in.SenderID,
		"receiver_id": in.ReceiverID,
	})

	socket.Session.Chat.Write(c)
}
