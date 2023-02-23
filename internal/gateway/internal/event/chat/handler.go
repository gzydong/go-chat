package chat

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/logger"
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
	h.handlers[entity.EventTalkKeyboard] = h.OnKeyboardMessage
	h.handlers[entity.EventTalkRead] = h.OnReadMessage

	// 聊天消息
	h.handlers["event.talk.text.message"] = h.OnTextMessage
	h.handlers["event.talk.image.message"] = h.OnImageMessage
	h.handlers["event.talk.file.message"] = h.OnFileMessage
	h.handlers["event.talk.code.message"] = h.OnCodeMessage
	h.handlers["event.talk.location.message"] = h.OnLocationMessage
	h.handlers["event.talk.vote.message"] = h.OnVoteMessage
}

func (h *Handler) Call(ctx context.Context, client im.IClient, event string, data []byte) {

	if h.handlers == nil {
		h.Init()
	}

	if call, ok := h.handlers[event]; ok {
		call(ctx, client, data)
	} else {
		logger.Warnf("Chat Event: [%s]未注册回调事件\n", event)
	}
}
