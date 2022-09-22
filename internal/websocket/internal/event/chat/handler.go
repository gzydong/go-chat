package chat

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/im"
	"go-chat/internal/service"
)

type Handler struct {
	redis         *redis.Client
	memberService *service.GroupMemberService
	handlers      map[string]func(ctx context.Context, client im.IClient, data []byte)
}

func NewHandler(redis *redis.Client, memberService *service.GroupMemberService) *Handler {
	return &Handler{redis: redis, memberService: memberService}
}

func (h *Handler) Init() {

	h.handlers = make(map[string]func(ctx context.Context, client im.IClient, data []byte))

	// 注册自定义绑定事件
	h.handlers[entity.EventTalkKeyboard] = h.OnKeyboard
	h.handlers[entity.EventTalkRead] = h.OnReadMessage
}

func (h *Handler) Call(ctx context.Context, client im.IClient, event string, data []byte) {

	if h.handlers == nil {
		h.Init()
	}

	if call, ok := h.handlers[event]; ok {
		call(ctx, client, data)
	} else {
		logrus.Warnf("Chat Event: [%s]未注册回调事件\n", event)
	}
}
