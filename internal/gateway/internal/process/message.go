package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/gateway/internal/consume"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/worker"
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
	// 订阅通道
	sub := m.redis.Subscribe(ctx, topic...)
	defer sub.Close()

	w := worker.NewWorker(10, 10)

	// 订阅 redis 消息
	for msg := range sub.Channel(redis.WithChannelHealthCheckInterval(30 * time.Second)) {
		w.Do(func() {
			var message *SubscribeContent
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				logger.Warnf("订阅消息格式错误 Err: %s \n", err.Error())
				return
			}

			// 触发回调方法
			consume.Call(message.Event, message.Data)
		})
	}

	w.Wait()
}
