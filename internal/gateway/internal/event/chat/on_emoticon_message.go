package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type EmoticonMessage struct {
	MsgId string                         `json:"msg_id"`
	Event string                         `json:"event"`
	Body  message.EmoticonMessageRequest `json:"body"`
}

// OnEmoticonMessage 表情包消息
func (h *Handler) OnEmoticonMessage(ctx context.Context, _ im.IClient, data []byte) {

	var m *EmoticonMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}

	fmt.Println("[OnEmoticonMessage] 新消息 ", string(data))
}
