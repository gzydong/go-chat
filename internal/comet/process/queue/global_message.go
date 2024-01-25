package queue

import (
	"context"
	"go-chat/internal/pkg/core/consumer"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/server"
)

var _ consumer.IConsumerHandle = (*GlobalMessage)(nil)

type GlobalMessage struct {
	Room *socket.RoomStorage
}

func (i *GlobalMessage) Topic() string {
	return "im.message.global"
}

func (i *GlobalMessage) Channel() string {
	return server.ID()
}

func (i *GlobalMessage) Touch() bool {
	return false
}

func (i *GlobalMessage) Do(ctx context.Context, message []byte, attempts uint16) error {
	return nil
}
