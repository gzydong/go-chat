package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type CodeMessage struct {
	MsgId string                     `json:"msg_id"`
	Event string                     `json:"event"`
	Body  message.CodeMessageRequest `json:"body"`
}

// OnCodeMessage 代码消息
func (h *Handler) OnCodeMessage(ctx context.Context, _ im.IClient, data []byte) {

	var m *CodeMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnCodeMessage Err: ", err)
		return
	}

	fmt.Println("[OnCodeMessage] 新消息 ", string(data))
}
