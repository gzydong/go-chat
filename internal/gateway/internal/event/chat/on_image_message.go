package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

type ImageMessage struct {
	MsgId   string                      `json:"msg_id"`
	Event   string                      `json:"event"`
	Content message.ImageMessageRequest `json:"content"`
}

// OnImageMessage 图片消息
func (h *Handler) OnImageMessage(ctx context.Context, _ socket.IClient, data []byte) {

	var m ImageMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnImageMessage Err: ", err)
		return
	}

	fmt.Println("[OnImageMessage] 新消息 ", string(data))
}
