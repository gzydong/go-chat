package process

import (
	"context"
	"errors"
	"github.com/nsqio/go-nsq"
	"go-chat/config"
	"go-chat/internal/comet/process/queue"
	"go-chat/internal/pkg/core/consumer"
	"log"
)

type QueueSubscribe struct {
	Config        *config.Config
	GlobalMessage *queue.GlobalMessage
	LocalMessage  *queue.LocalMessage
	RoomControl   *queue.RoomControl
}

func (m *QueueSubscribe) Setup(ctx context.Context) error {

	c := consumer.NewConsumer(m.Config.Nsq.Addr, nsq.NewConfig())

	c.Register("default", m.GlobalMessage)
	c.Register("default", m.RoomControl)
	c.Register("default", m.LocalMessage)

	if err := c.Start(ctx, "default"); err != nil {
		log.Fatal(err)
	}

	return errors.New("not implement")
}
