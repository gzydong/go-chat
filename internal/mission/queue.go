package mission

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
	"go-chat/internal/mission/queue"
)

type QueueProvider struct {
	Consumers *queue.Consumers
	Redis     *redis.Client
}

func Queue(ctx *cli.Context, app *QueueProvider) error {
	topics := []string{"im.user.login"}

	sub := app.Redis.Subscribe(ctx.Context, topics...)

	// nolint
	defer sub.Close()

	log.Printf("subscribed to topics: %v", topics)

	for data := range sub.Channel(redis.WithChannelHealthCheckInterval(10 * time.Second)) {
		switch data.Channel {
		case "im.user.login":
			_ = app.Consumers.UserLoginConsumer.Do(context.Background(), []byte(data.Payload), 1)
		}
	}

	return nil
}
