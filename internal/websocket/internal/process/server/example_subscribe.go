package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/worker"
	"go-chat/internal/websocket/internal/process/consume"
)

type ExampleSubscribe struct {
	redis   *redis.Client
	config  *config.Config
	consume *consume.ExampleSubscribe
}

func NewExampleSubscribe(redis *redis.Client, config *config.Config, consume *consume.ExampleSubscribe) *ExampleSubscribe {
	return &ExampleSubscribe{redis: redis, config: config, consume: consume}
}

func (e *ExampleSubscribe) Setup(ctx context.Context) error {

	log.Println("ExampleSubscribe Setup")

	gateway := fmt.Sprintf(entity.ImTopicExamplePrivate, e.config.ServerId())

	channels := []string{
		entity.ImTopicExample, // 全局通道
		gateway,               // 私有通道
	}

	// 订阅通道
	sub := e.redis.Subscribe(ctx, channels...)
	defer sub.Close()

	go e.subscribe(sub, func(value *redis.Message) {
		fmt.Println(value.Payload)
	})

	<-ctx.Done()

	return nil
}

func (*ExampleSubscribe) subscribe(sub *redis.PubSub, fn func(value *redis.Message)) {
	go func() {
		w := worker.NewWorker(10, 10)

		// 订阅 redis 消息
		for msg := range sub.Channel(redis.WithChannelHealthCheckInterval(30 * time.Second)) {
			w.Do(func() {
				fn(msg)
			})
		}

		w.Wait()
	}()
}
