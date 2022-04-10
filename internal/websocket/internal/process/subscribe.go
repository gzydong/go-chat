package process

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/pool"
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

func (w *WsSubscribe) Handle(ctx context.Context) error {
	gateway := fmt.Sprintf(entity.IMGatewayPrivate, w.conf.ServerId())

	channels := []string{
		entity.IMGatewayAll, // 全局通道
		gateway,             // 私有通道
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	defer sub.Close()

	go func() {
		work := pool.NewWorkerPool(5) // 设置协程并发处理数

		for msg := range sub.Channel(redis.WithChannelHealthCheckInterval(30 * time.Second)) {
			consume := func(value *redis.Message) {
				switch value.Channel {

				// 私有通道及全局广播通道
				case gateway, entity.IMGatewayAll:
					var message *SubscribeContent
					if err := json.Unmarshal([]byte(value.Payload), &message); err == nil {
						w.consume.Handle(message.Event, message.Data)
					}
				}
			}

			work.Add(func() { consume(msg) })
		}

		work.Wait()
	}()

	<-ctx.Done()

	return nil
}
