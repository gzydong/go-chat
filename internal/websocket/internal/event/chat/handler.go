package chat

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go-chat/internal/entity"
)

type Handler struct {
	redis    *redis.Client
	handlers map[string]func(data string)
}

func NewHandler(redis *redis.Client) *Handler {
	return &Handler{redis: redis}
}

func (h *Handler) Init() {

	h.handlers = make(map[string]func(data string))

	// 注册自定义绑定事件
	h.handlers[entity.EventTalkKeyboard] = h.OnKeyboard
	h.handlers[entity.EventTalkRead] = h.OnReadMessage
}

func (h *Handler) Call(ctx context.Context, event string, data string) {

	if h.handlers == nil {
		h.Init()
	}

	if call, ok := h.handlers[event]; ok {
		call(data)
	} else {
		logrus.Warnf("Chat Event: [%s]未注册回调事件\n", event)
	}
}
