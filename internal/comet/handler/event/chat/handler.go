package chat

import (
	"context"
	"go-chat/internal/business"
	"log"

	"github.com/redis/go-redis/v9"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type handle func(ctx context.Context, client socket.IClient, data []byte)

var handlers map[string]handle

type Handler struct {
	Redis         *redis.Client
	Source        *repo.Source
	MemberService service.IGroupMemberService
	PushMessage   *business.PushMessage
}

func (h *Handler) init() {
	handlers = make(map[string]handle)
	// 注册自定义绑定事件
	handlers["im.message.keyboard"] = h.onKeyboardMessage
}

func (h *Handler) Call(ctx context.Context, client socket.IClient, event string, data []byte) {

	if handlers == nil {
		h.init()
	}

	if call, ok := handlers[event]; ok {
		call(ctx, client, data)
	} else {
		log.Printf("Chat Event: [%s]未注册回调事件\n", event)
	}
}
