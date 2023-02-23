package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type ImageMessage struct {
	MsgId string                      `json:"msg_id"`
	Event string                      `json:"event"`
	Body  message.ImageMessageRequest `json:"body"`
}

// OnImageMessage 图片消息
func (h *Handler) OnImageMessage(ctx context.Context, _ im.IClient, data []byte) {

	var m *ImageMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}

	fmt.Println("[OnImageMessage] 新消息 ", string(data))
}
