package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type TextMessage struct {
	MsgId string                     `json:"msg_id"`
	Event string                     `json:"event"`
	Body  message.TextMessageRequest `json:"body"`
}

// OnTextMessage 文本消息
func (h *Handler) OnTextMessage(ctx context.Context, client im.IClient, data []byte) {

	var m *TextMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}

	fmt.Println("[TextMessage] 新消息 ", string(data))
}
