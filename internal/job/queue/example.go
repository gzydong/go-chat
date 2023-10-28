package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

type ExampleQueue struct {
}

func (e *ExampleQueue) Handle(ctx context.Context) error {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second

	c, err := nsq.NewConsumer("test", "test_channel", config)
	if err != nil {
		return err
	}

	c.AddHandler(e)

	return c.ConnectToNSQD("127.0.0.1:4150")
}

// HandleMessage 是需要实现的处理消息的方法
func (e *ExampleQueue) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("%s recv from %v, msg:%v\n", "ExampleQueue", msg.NSQDAddress, string(msg.Body))
	return nil
}
