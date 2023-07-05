package example

import (
	"context"
	"log"
)

type Handler struct {
	handlers map[string]func(ctx context.Context, data []byte)
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) init() {
	h.handlers = make(map[string]func(ctx context.Context, data []byte))
}

func (h *Handler) Call(ctx context.Context, event string, data []byte) {
	if h.handlers == nil {
		h.init()
	}

	if call, ok := h.handlers[event]; ok {
		call(ctx, data)
	} else {
		log.Printf("consume chat event: [%s]未注册回调事件\n", event)
	}
}
