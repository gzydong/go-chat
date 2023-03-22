package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

type FileMessage struct {
	MsgId   string                      `json:"msg_id"`
	Event   string                      `json:"event"`
	Content message.ImageMessageRequest `json:"content"`
}

// OnFileMessage 文本消息
func (h *Handler) OnFileMessage(ctx context.Context, _ socket.IClient, data []byte) {

	var m *FileMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnFileMessage Err: ", err)
		return
	}

	fmt.Println("[OnFileMessage] 新消息 ", string(data))
}
