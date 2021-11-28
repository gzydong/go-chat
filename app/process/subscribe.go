package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/cache"
	"go-chat/app/entity"
	"go-chat/app/process/handle"
	"go-chat/app/service"
	"go-chat/config"
)

type MessagePayload struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type WsSubscribe struct {
	rds                *redis.Client
	conf               *config.Config
	talkRecordsService *service.TalkRecordsService
	ws                 *cache.WsClientSession
	room               *cache.GroupRoom
	contactService     *service.ContactService
	consume            *handle.SubscribeConsume
}

func NewWsSubscribe(rds *redis.Client, conf *config.Config, talkRecordsService *service.TalkRecordsService, ws *cache.WsClientSession, room *cache.GroupRoom, contactService *service.ContactService, consume *handle.SubscribeConsume) *WsSubscribe {
	return &WsSubscribe{rds: rds, conf: conf, talkRecordsService: talkRecordsService, ws: ws, room: room, contactService: contactService, consume: consume}
}

type SubscribeContent struct {
	Event string `json:"event_name"`
	Data  string `json:"data"`
}

func (w *WsSubscribe) Handle(ctx context.Context) error {
	gateway := fmt.Sprintf(entity.SubscribeWsGatewayPrivate, w.conf.GetSid())

	channels := []string{
		entity.SubscribeWsGatewayAll, // 全局通道
		gateway,                      // 私有通道
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	defer sub.Close()

	go func() {
		for msg := range sub.Channel() {
			fmt.Printf("消息订阅 : channel=%s message=%s\n", msg.Channel, msg.Payload)

			switch msg.Channel {
			case gateway, entity.SubscribeWsGatewayAll:
				var message *SubscribeContent

				if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
					continue
				}

				w.consume.Handle(message.Event, message.Data)
			}
		}
	}()

	<-ctx.Done()

	return nil
}
