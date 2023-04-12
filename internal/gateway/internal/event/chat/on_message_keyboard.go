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
	Event   string `json:"event"`
	Content struct {
		SenderID   int `json:"sender_id"`
		ReceiverID int `json:"receiver_id"`
	} `json:"content"`
}

// onKeyboardMessage 键盘输入事件
func (h *Handler) onKeyboardMessage(ctx context.Context, _ socket.IClient, data []byte) {

	var m KeyboardMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat onKeyboardMessage Err: ", err)
		return
	}

	h.redis.Publish(ctx, entity.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": entity.SubEventImMessageKeyboard,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   m.Content.SenderID,
			"receiver_id": m.Content.ReceiverID,
		}),
	}))
}
