package chat

import (
	"context"
	"encoding/json"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/websocket/internal/dto"
)

// OnKeyboard 键盘输入事件
func (h *Handler) OnKeyboard(data string) {
	var m *dto.KeyboardMessage

	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return
	}

	h.redis.Publish(context.Background(), entity.ImTopicDefault, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventTalkKeyboard,
		"data": jsonutil.Encode(entity.MapStrAny{
			"sender_id":   m.Data.SenderID,
			"receiver_id": m.Data.ReceiverID,
		}),
	}))
}
