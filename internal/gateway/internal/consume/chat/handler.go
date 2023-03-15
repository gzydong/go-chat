package chat

import (
	"context"
	"log"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/repository/cache"
	"go-chat/internal/service"
)

type Handler struct {
	handlers map[string]func(data []byte)

	config         *config.Config
	clientStorage  *cache.ClientStorage
	roomStorage    *cache.RoomStorage
	recordsService *service.TalkRecordsService
	contactService *service.ContactService
}

func NewHandler(config *config.Config, clientStorage *cache.ClientStorage, roomStorage *cache.RoomStorage, recordsService *service.TalkRecordsService, contactService *service.ContactService) *Handler {
	return &Handler{config: config, clientStorage: clientStorage, roomStorage: roomStorage, recordsService: recordsService, contactService: contactService}
}

func (h *Handler) init() {
	h.handlers = make(map[string]func(data []byte))

	h.handlers[entity.EventTalk] = h.onConsumeTalk
	h.handlers[entity.EventTalkKeyboard] = h.onConsumeTalkKeyboard
	h.handlers[entity.EventOnlineStatus] = h.onConsumeLogin
	h.handlers[entity.EventTalkRevoke] = h.onConsumeTalkRevoke
	h.handlers[entity.EventTalkJoinGroup] = h.onConsumeTalkJoinGroup
	h.handlers[entity.EventContactApply] = h.onConsumeContactApply
	h.handlers[entity.EventTalkRead] = h.onConsumeTalkRead
}

func (h *Handler) Call(ctx context.Context, client socket.IClient, event string, data []byte) {

	if h.handlers == nil {
		h.init()
	}

	if call, ok := h.handlers[event]; ok {
		call(data)
	} else {
		log.Printf("consume chat event: [%s]未注册回调事件\n", event)
	}
}
