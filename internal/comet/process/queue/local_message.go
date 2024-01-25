package queue

import (
	"context"
	"fmt"
	"go-chat/internal/pkg/core/consumer"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/pkg/server"
)

var _ consumer.IConsumerHandle = (*LocalMessage)(nil)

type LocalMessage struct {
	Room *socket.RoomStorage
}

func (i *LocalMessage) Topic() string {
	return fmt.Sprintf("im.message.local.%s", server.ID())
}

func (i *LocalMessage) Channel() string {
	return "default"
}

func (i *LocalMessage) Touch() bool {
	return false
}

func (i *LocalMessage) Do(ctx context.Context, message []byte, attempts uint16) error {
	return nil
}
