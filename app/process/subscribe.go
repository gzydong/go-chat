package process

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-chat/app/pkg/im"
)

type MessagePayload struct {
	EventName string `json:"event_name"`
	Data      string `json:"data"`
}

type WsSubscribe struct {
	rds *redis.Client
}

func NewWsSubscribe(rds *redis.Client) *WsSubscribe {
	return &WsSubscribe{rds: rds}
}

func (w *WsSubscribe) Handle(ctx context.Context) error {
	sub := w.rds.Subscribe(ctx, "chat")
	defer sub.Close()

	go func() {
		for msg := range sub.Channel() {
			body := im.NewSenderContent()
			body.SetBroadcast(true)
			body.SetMessage(&im.Message{
				Event:   "talk",
				Content: msg.Payload,
			})

			im.GroupManage.DefaultChannel.PushSendChannel(body)

			fmt.Printf("channel=%s message=%s\n", msg.Channel, msg.Payload)
		}
	}()

	<-ctx.Done()

	return nil
}
