package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type FileMessage struct {
	MsgId string                      `json:"msg_id"`
	Event string                      `json:"event"`
	Body  message.ImageMessageRequest `json:"body"`
}

// OnFileMessage 文本消息
func (h *Handler) OnFileMessage(ctx context.Context, _ im.IClient, data []byte) {

	var m *FileMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}

	fmt.Println("[OnFileMessage] 新消息 ", string(data))
}
