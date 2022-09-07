package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/worker"
	"go-chat/internal/websocket/internal/process/consume"
)

type SubscribeContent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type DefaultSubscribe struct {
	redis   *redis.Client
	config  *config.Config
	consume *consume.DefaultSubscribe
}

func NewDefaultSubscribe(redis *redis.Client, config *config.Config, consume *consume.DefaultSubscribe) *DefaultSubscribe {
	return &DefaultSubscribe{redis: redis, config: config, consume: consume}
}

func (d *DefaultSubscribe) Setup(ctx context.Context) error {

	log.Println("DefaultSubscribe Setup")

	gateway := fmt.Sprintf(entity.ImTopicDefaultPrivate, d.config.ServerId())

	channels := []string{
		entity.ImTopicDefault, // 全局通道
		gateway,               // 私有通道
	}

	// 订阅通道
	sub := d.redis.Subscribe(ctx, channels...)
	defer sub.Close()

	// 注册处理事件
	d.consume.RegisterEvent()
	go d.subscribe(sub, func(value *redis.Message) {
		switch value.Channel {

		// 私有通道及全局广播通道
		case gateway, entity.ImTopicDefault:
			var message *SubscribeContent
			if err := json.Unmarshal([]byte(value.Payload), &message); err == nil {
				d.consume.Handle(message.Event, message.Data)
			} else {
				logger.Warnf("订阅消息格式错误 Err: %s \n", err.Error())
			}
		}
	})

	<-ctx.Done()

	return nil
}

func (*DefaultSubscribe) subscribe(sub *redis.PubSub, fn func(value *redis.Message)) {
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
