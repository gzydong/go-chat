package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/entity"
	"go-chat/app/pkg/pool"
	"go-chat/app/process/handle"
	"go-chat/config"
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
	gateway := fmt.Sprintf(entity.SubscribeWsGatewayPrivate, w.conf.GetSid())

	channels := []string{
		entity.SubscribeWsGatewayAll, // 全局通道
		gateway,                      // 私有通道
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	_, _ = sub.ReceiveTimeout(ctx, 0)

	defer sub.Close()

	go func() {
		work := pool.NewWorkerPool(5) // 设置协程并发处理数

		for msg := range sub.Channel() {
			// fmt.Printf("消息订阅 : channel=%s message=%s\n", msg.channel, msg.Payload)

			consume := func(value *redis.Message) {
				switch value.Channel {

				// 私有通道及全局广播通道
				case gateway, entity.SubscribeWsGatewayAll:
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
