package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sourcegraph/conc/pool"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/gateway/internal/consume"
)

type MessageSubscribe struct {
	config         *config.Config
	redis          *redis.Client
	defaultConsume *consume.ChatSubscribe
	exampleConsume *consume.ExampleSubscribe
}

func NewMessageSubscribe(config *config.Config, redis *redis.Client, defaultConsume *consume.ChatSubscribe, exampleConsume *consume.ExampleSubscribe) *MessageSubscribe {
	return &MessageSubscribe{config: config, redis: redis, defaultConsume: defaultConsume, exampleConsume: exampleConsume}
}

type IConsume interface {
	Call(event string, data string)
}

func (m *MessageSubscribe) Setup(ctx context.Context) error {

	log.Println("Start MessageSubscribe")

	go m.subscribe(ctx, []string{entity.ImTopicChat, fmt.Sprintf(entity.ImTopicChatPrivate, m.config.ServerId())}, m.defaultConsume)

	go m.subscribe(ctx, []string{entity.ImTopicExample, fmt.Sprintf(entity.ImTopicExamplePrivate, m.config.ServerId())}, m.exampleConsume)

	<-ctx.Done()

	return nil
}

type SubscribeContent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func (m *MessageSubscribe) subscribe(ctx context.Context, topic []string, consume IConsume) {

	sub := m.redis.Subscribe(ctx, topic...)
	defer sub.Close()

	worker := pool.New().WithMaxGoroutines(24)

	// 订阅 redis 消息
	for msg := range sub.Channel(redis.WithChannelHealthCheckInterval(15 * time.Second)) {
		worker.Go(func() {
			var message *SubscribeContent
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				log.Println("SubscribeContent Err: ", err.Error())
				return
			}

			consume.Call(message.Event, message.Data)
		})
	}

	worker.Wait()
}
