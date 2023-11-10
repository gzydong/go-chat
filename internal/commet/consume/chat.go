package consume

import (
	"context"

	"go-chat/internal/commet/consume/chat"
)

type ChatSubscribe struct {
	handler *chat.Handler
}

func NewChatSubscribe(handel *chat.Handler) *ChatSubscribe {
	return &ChatSubscribe{handler: handel}
}

// Call 触发回调事件
func (s *ChatSubscribe) Call(event string, data []byte) {
	s.handler.Call(context.TODO(), event, data)
}
