package example

import (
	"context"

	"github.com/sirupsen/logrus"
	"go-chat/internal/pkg/im"
)

type Handler struct {
	handlers map[string]func(ctx context.Context, client im.IClient, data []byte)
}

func (h *Handler) Init() {

	h.handlers = make(map[string]func(ctx context.Context, client im.IClient, data []byte))

	// 注册自定义绑定事件
	// h.handlers[entity.EventTalkKeyboard] = h.OnKeyboard
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
