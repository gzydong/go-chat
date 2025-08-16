package consume

import (
	"context"
	"encoding/json"
	"log/slog"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
)

// 键盘输入事件消息
func (h *Handler) onConsumeTalkKeyboard(ctx context.Context, body []byte) {
	var in entity.SubEventImMessageKeyboardPayload
	if err := json.Unmarshal(body, &in); err != nil {
		logger.Errorf("[ChatSubscribe] onConsumeTalkKeyboard Unmarshal err: %s", err.Error())
		return
	}

	data := Message(entity.PushEventImMessageKeyboard, entity.ImMessageKeyboardPayload{
		FromId:   in.ToFromId,
		ToFromId: in.ToFromId,
	})

	for _, session := range h.serv.SessionManager().GetSessions(int64(in.ToFromId)) {
		if err := session.Write(data); err != nil {
			slog.Error("session write message error", "error", err)
		}
	}
}
