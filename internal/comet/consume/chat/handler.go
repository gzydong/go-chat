package chat

import (
	"context"
	"log"

	"go-chat/internal/pkg/core/socket"

	"go-chat/config"
	"go-chat/internal/entity"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

var handlers map[string]func(ctx context.Context, data []byte)

type Handler struct {
	Config               *config.Config
	OrganizeRepo         *repo.Organize
	UserRepo             *repo.Users
	Source               *repo.Source
	TalkRecordsService   service.ITalkRecordService
	ContactService       service.IContactService
	ClientConnectService service.IClientConnectService
	RoomStorage          *socket.RoomStorage
}

func (h *Handler) init() {
	handlers = make(map[string]func(ctx context.Context, data []byte))

	handlers[entity.SubEventImMessage] = h.onConsumeTalk
	handlers[entity.SubEventImMessageKeyboard] = h.onConsumeTalkKeyboard
	handlers[entity.SubEventImMessageRevoke] = h.onConsumeTalkRevoke
	handlers[entity.SubEventContactStatus] = h.onConsumeContactStatus
	handlers[entity.SubEventContactApply] = h.onConsumeContactApply
	handlers[entity.SubEventGroupJoin] = h.onConsumeGroupJoin
	handlers[entity.SubEventGroupApply] = h.onConsumeGroupApply
}

func (h *Handler) Call(ctx context.Context, event string, data []byte) {
	if handlers == nil {
		h.init()
	}

	if call, ok := handlers[event]; ok {
		call(ctx, data)
	} else {
		log.Printf("consume chat event: [%s]未注册回调事件\n", event)
	}
}
