package chat

import (
	"context"
	"encoding/json"
	"log"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/jsonutil"
)

type KeyboardMessage struct {
	Event string `json:"event"`
	Data  struct {
		SenderID   int `json:"sender_id"`
		ReceiverID int `json:"receiver_id"`
	} `json:"data"`
}

// OnKeyboardMessage 键盘输入事件
func (h *Handler) OnKeyboardMessage(ctx context.Context, _ socket.IClient, data []byte) {

	var m *KeyboardMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnKeyboardMessage Err: ", err)
		return
	}

	h.redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(entity.MapStrAny{
		"event": entity.EventTalkKeyboard,
		"data": jsonutil.Encode(entity.MapStrAny{
			"sender_id":   m.Data.SenderID,
			"receiver_id": m.Data.ReceiverID,
		}),
	}))
}
