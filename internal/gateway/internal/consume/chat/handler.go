package chat

import (
	"context"
	"log"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
	"go-chat/internal/service"
)

type Handler struct {
	handlers map[string]func(ctx context.Context, data []byte)

	config         *config.Config
	clientStorage  *cache.ClientStorage
	roomStorage    *cache.RoomStorage
	recordsService *service.TalkRecordsService
	contactService *service.ContactService
	organize       *organize.Organize
	source         *repo.Source
}

func NewHandler(config *config.Config, clientStorage *cache.ClientStorage, roomStorage *cache.RoomStorage, recordsService *service.TalkRecordsService, contactService *service.ContactService, organize *organize.Organize, source *repo.Source) *Handler {
	return &Handler{config: config, clientStorage: clientStorage, roomStorage: roomStorage, recordsService: recordsService, contactService: contactService, organize: organize, source: source}
}

func (h *Handler) init() {
	h.handlers = make(map[string]func(ctx context.Context, data []byte))

	h.handlers[entity.SubEventImMessage] = h.onConsumeTalk
	h.handlers[entity.SubEventImMessageKeyboard] = h.onConsumeTalkKeyboard
	h.handlers[entity.SubEventImMessageRead] = h.onConsumeTalkRead
	h.handlers[entity.SubEventImMessageRevoke] = h.onConsumeTalkRevoke
	h.handlers[entity.SubEventContactStatus] = h.onConsumeContactStatus
	h.handlers[entity.SubEventContactApply] = h.onConsumeContactApply
	h.handlers[entity.SubEventGroupJoin] = h.onConsumeGroupJoin
	h.handlers[entity.SubEventGroupApply] = h.onConsumeGroupApply
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
