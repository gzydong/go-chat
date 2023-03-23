package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

type CodeMessage struct {
	MsgId   string                     `json:"msg_id"`
	Event   string                     `json:"event"`
	Content message.CodeMessageRequest `json:"content"`
}

// OnCodeMessage 代码消息
func (h *Handler) OnCodeMessage(ctx context.Context, _ socket.IClient, data []byte) {

	var m CodeMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnCodeMessage Err: ", err)
		return
	}

	fmt.Println("[OnCodeMessage] 新消息 ", string(data))
}
