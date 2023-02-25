package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type LocationMessageMessage struct {
	MsgId string                         `json:"msg_id"`
	Event string                         `json:"event"`
	Body  message.LocationMessageRequest `json:"body"`
}

// OnLocationMessage 位置消息
func (h *Handler) OnLocationMessage(ctx context.Context, _ im.IClient, data []byte) {

	var m *LocationMessageMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnLocationMessage Err: ", err)
		return
	}

	fmt.Println("[OnLocationMessage] 新消息 ", string(data))
}
