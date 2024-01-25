package chat

import (
	"context"
	"encoding/json"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/server"
)

// 键盘输入事件消息
func (h *Handler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {
	var in entity.SubEventImMessageKeyboardPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: %s", err.Error())
		return
	}

	ids, _ := h.ClientConnectService.GetUidFromClientIds(ctx, server.ID(), socket.Session.Chat.Name(), in.ToFromId)
	if len(ids) == 0 {
		return
	}

	c := socket.NewSenderContent()
	c.SetReceive(ids...)
	c.SetMessage(entity.PushEventImMessageKeyboard, entity.ImMessageKeyboardPayload{
		FromId:   in.ToFromId,
		ToFromId: in.ToFromId,
	})

	socket.Session.Chat.Write(c)
}
