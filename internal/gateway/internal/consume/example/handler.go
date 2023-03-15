package example

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"go-chat/internal/pkg/ichat/socket"
)

type Handler struct {
	redis    *redis.Client
	handlers map[string]func(ctx context.Context, client socket.IClient, data []byte)
}

func NewHandler(redis *redis.Client) *Handler {
	return &Handler{redis: redis}
}

func (h *Handler) Init() {
	h.handlers = make(map[string]func(ctx context.Context, client socket.IClient, data []byte))
}

func (h *Handler) Call(ctx context.Context, client socket.IClient, event string, data []byte) {

	if h.handlers == nil {
		h.Init()
	}

	if call, ok := h.handlers[event]; ok {
		call(ctx, client, data)
	} else {
		log.Printf("consume example event: [%s]未注册回调事件\n", event)
	}
}
