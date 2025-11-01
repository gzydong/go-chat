package comet

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gzydong/go-chat/internal/comet/consume"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/longnet"
	"github.com/gzydong/go-chat/internal/pkg/server"
	"github.com/gzydong/go-chat/internal/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/pool"
)

var _ longnet.IProcess = &Subscribe{}

type Subscribe struct {
	Redis   *redis.Client
	Handler *consume.Handler
}

func (m *Subscribe) Start(ctx context.Context, serv longnet.IServer) error {
	m.Handler.SetServ(serv)

	sub := m.Redis.Subscribe(ctx, []string{entity.ImChannelChat, entity.ImTopicChat, fmt.Sprintf(entity.ImTopicChatPrivate, server.ID())}...)
	defer func() {
		_ = sub.Close()
	}()

	worker := pool.New().WithMaxGoroutines(10)
	defer worker.Wait()

	for {
		select {
		case <-ctx.Done():
			return nil
		case data := <-sub.Channel():
			worker.Go(func() {
				m.handle(data)
			})
		}
	}
}

func (m *Subscribe) handle(data *redis.Message) {
	var in entity.SubscribeMessage
	if err := json.Unmarshal([]byte(data.Payload), &in); err != nil {
		slog.Error("[payload] subscribe content unmarshal Err: ", "error", err.Error())
		return
	}

	defer func() {
		if err := recover(); err != nil {
			slog.Error("message subscribe call err", "panic", utils.PanicTrace(err))
		}
	}()

	m.Handler.Call(context.Background(), in.Event, []byte(in.Payload))
}
