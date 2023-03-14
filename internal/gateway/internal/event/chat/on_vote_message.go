package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/ichat/socket"
)

type VoteMessage struct {
	MsgId string                     `json:"msg_id"`
	Event string                     `json:"event"`
	Body  message.VoteMessageRequest `json:"body"`
}

// OnVoteMessage 文本消息
func (h *Handler) OnVoteMessage(ctx context.Context, _ socket.IClient, data []byte) {

	var m *VoteMessage
	if err := json.Unmarshal(data, &m); err != nil {
		log.Println("Chat OnVoteMessage Err: ", err)
		return
	}

	fmt.Println("[OnVoteMessage] 新消息 ", string(data))
}
