package example

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	handlers map[string]func(data string)
}

func (h *Handler) Init() {

	h.handlers = make(map[string]func(data string))

	// 注册自定义绑定事件
	// h.handlers[entity.EventTalkKeyboard] = h.OnKeyboard
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
