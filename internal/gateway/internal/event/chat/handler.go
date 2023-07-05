package chat

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/service"
)

type handle func(ctx context.Context, client socket.IClient, data []byte)

type Handler struct {
	redis         *redis.Client
	memberService *service.GroupMemberService
	handlers      map[string]func(ctx context.Context, client socket.IClient, data []byte)
	message       *service.MessageService
}

func NewHandler(redis *redis.Client, memberService *service.GroupMemberService, message *service.MessageService) *Handler {
	return &Handler{redis: redis, memberService: memberService, message: message}
}

func (h *Handler) init() {

	h.handlers = make(map[string]func(ctx context.Context, client socket.IClient, data []byte))

	// 注册自定义绑定事件
	h.handlers["im.message.publish"] = h.onPublish
	h.handlers["im.message.revoke"] = h.onRevokeMessage
	h.handlers["im.message.delete"] = h.onDeleteMessage
	h.handlers["im.message.read"] = h.onReadMessage
	h.handlers["im.message.keyboard"] = h.onKeyboardMessage
}

func (h *Handler) Call(ctx context.Context, client socket.IClient, event string, data []byte) {

	if h.handlers == nil {
		h.init()
	}

	if call, ok := h.handlers[event]; ok {
		call(ctx, client, data)
	} else {
		log.Printf("Chat Event: [%s]未注册回调事件\n", event)
	}
}
