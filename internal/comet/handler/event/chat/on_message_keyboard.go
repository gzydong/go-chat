package chat

import (
	"context"
	"encoding/json"
	"log"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/jsonutil"
)

type KeyboardMessage struct {
	Event   string `json:"event"`
	Payload struct {
		ToFromId int `json:"to_from_id"`
	} `json:"payload"`
}

// onKeyboardMessage 键盘输入事件
func (h *Handler) onKeyboardMessage(ctx context.Context, c socket.IClient, data []byte) {
	var in KeyboardMessage
	if err := json.Unmarshal(data, &in); err != nil {
		log.Println("Chat onKeyboardMessage Err: ", err)
		return
	}

	_ = h.PushMessage.Push(ctx, entity.ImTopicChat, &entity.SubscribeMessage{
		Event: entity.SubEventImMessageKeyboard,
		Payload: jsonutil.Encode(entity.SubEventImMessageKeyboardPayload{
			FromId:   c.Uid(),
			ToFromId: in.Payload.ToFromId,
		}),
	})
}
