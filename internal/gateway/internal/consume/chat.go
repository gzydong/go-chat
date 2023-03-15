package consume

import (
	"context"

	"go-chat/internal/gateway/internal/consume/chat"
)

type ChatSubscribe struct {
	handler *chat.Handler
}

func NewChatSubscribe(handel *chat.Handler) *ChatSubscribe {
	return &ChatSubscribe{handler: handel}
}

// Call 触发回调事件
func (s *ChatSubscribe) Call(event string, data string) {
	s.handler.Call(context.Background(), event, []byte(data))
}
