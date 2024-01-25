package process

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
	"go-chat/config"
	"go-chat/internal/commet/consume"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/utils"
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
	Call(event string, data []byte)
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

	worker := pool.New().WithMaxGoroutines(10)

	for data := range sub.Channel(redis.WithChannelHealthCheckInterval(10 * time.Second)) {
		m.handle(worker, data, consume)
	}

	worker.Wait()
}

func (m *MessageSubscribe) handle(worker *pool.Pool, data *redis.Message, consume IConsume) {
	worker.Go(func() {
		var in SubscribeContent
		if err := json.Unmarshal([]byte(data.Payload), &in); err != nil {
			log.Println("SubscribeContent Unmarshal Err: ", err.Error())
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Println("MessageSubscribe Call Err: ", utils.PanicTrace(err))
			}
		}()

		consume.Call(in.Event, []byte(in.Data))
	})
}
