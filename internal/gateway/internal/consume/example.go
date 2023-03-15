package consume

import (
	"context"

	"go-chat/internal/gateway/internal/consume/example"
)

type ExampleSubscribe struct {
	handler *example.Handler
}

func NewExampleSubscribe(handler *example.Handler) *ExampleSubscribe {
	return &ExampleSubscribe{handler: handler}
}

// Call 触发回调事件
func (s *ExampleSubscribe) Call(event string, data string) {
	s.handler.Call(context.Background(), event, []byte(data))
}
