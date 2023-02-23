package chat

import (
	"context"
	"encoding/json"
	"fmt"

	"go-chat/api/pb/message/v1"
	"go-chat/internal/pkg/im"
)

type VoteMessage struct {
	MsgId string                     `json:"msg_id"`
	Event string                     `json:"event"`
	Body  message.VoteMessageRequest `json:"body"`
}

// OnVoteMessage 文本消息
func (h *Handler) OnVoteMessage(ctx context.Context, _ im.IClient, data []byte) {

	var m *VoteMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}

	fmt.Println("[OnVoteMessage] 新消息 ", string(data))
}
