package example

import (
	"context"
	"fmt"

	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/pkg/logger"
)

type Handler struct {
	handlers map[string]func(ctx context.Context, client socket.IClient, data []byte)
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) init() {

	h.handlers = make(map[string]func(ctx context.Context, client socket.IClient, data []byte))

	// 注册自定义绑定事件
}

func (h *Handler) Call(ctx context.Context, client socket.IClient, event string, data []byte) {

	if h.handlers == nil {
		h.init()
	}

	if call, ok := h.handlers[event]; ok {
		call(ctx, client, data)
	} else {
		logger.Std().Warn(fmt.Sprintf("Chat Event: [%s]未注册回调事件\n", event))
	}
}
