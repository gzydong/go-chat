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
	"go-chat/internal/websocket/internal/process/handle"
)

type SubscribeContent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type WsSubscribe struct {
	rds     *redis.Client
	conf    *config.Config
	consume *handle.SubscribeConsume
}

func NewWsSubscribe(rds *redis.Client, conf *config.Config, consume *handle.SubscribeConsume) *WsSubscribe {
	return &WsSubscribe{rds: rds, conf: conf, consume: consume}
}

func (w *WsSubscribe) Setup(ctx context.Context) error {

	log.Println("WsSubscribe Setup")

	gateway := fmt.Sprintf(entity.IMGatewayPrivate, w.conf.ServerId())

	channels := []string{
		entity.IMGatewayAll, // 全局通道
		gateway,             // 私有通道
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	// 关闭订阅
	defer sub.Close()

	consume := func(value *redis.Message) {

		switch value.Channel {

		// 私有通道及全局广播通道
		case gateway, entity.IMGatewayAll:
			var message *SubscribeContent
			if err := json.Unmarshal([]byte(value.Payload), &message); err == nil {
				w.consume.Handle(message.Event, message.Data)
			} else {
				logger.Warnf("订阅消息格式错误 Err: %s \n", err.Error())
			}
		}
	}

	go func() {
		w := worker.NewWorker(10, 10)

		// 订阅 redis 消息
		for msg := range sub.Channel(redis.WithChannelHealthCheckInterval(30 * time.Second)) {
			w.Do(func() {
				consume(msg)
			})
		}

		w.Wait()
	}()

	<-ctx.Done()

	return nil
}
