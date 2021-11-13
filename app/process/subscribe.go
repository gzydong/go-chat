package process

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/pkg/im"
	"go-chat/config"
)

type MessagePayload struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type WsSubscribe struct {
	rds  *redis.Client
	conf *config.Config
}

func NewWsSubscribe(rds *redis.Client, conf *config.Config) *WsSubscribe {
	return &WsSubscribe{rds: rds, conf: conf}
}

func (w *WsSubscribe) Handle(ctx context.Context) error {
	channels := []string{
		"ws:all",                              // 全局通道
		fmt.Sprintf("ws:%s", w.conf.GetSid()), // 私有通道
	}

	// 订阅通道
	sub := w.rds.Subscribe(ctx, channels...)

	defer sub.Close()

	go func() {
		for msg := range sub.Channel() {
			body := im.NewSenderContent()
			body.SetBroadcast(true)
			body.SetMessage(&im.Message{
				Event:   "talk",
				Content: msg.Payload,
			})

			im.Session.DefaultChannel.PushSendChannel(body)

			fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
		}
	}()

	<-ctx.Done()

	return nil
}
