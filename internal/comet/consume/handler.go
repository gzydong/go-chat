package consume

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gzydong/go-chat/internal/pkg/longnet"

	"github.com/gzydong/go-chat/config"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
)

var handlers map[string]func(ctx context.Context, data []byte)

type Handler struct {
	Config             *config.Config
	OrganizeRepo       *repo.Organize
	UserRepo           *repo.Users
	Source             *repo.Source
	TalkRecordsService service.ITalkRecordService
	ContactService     service.IContactService
	serv               longnet.IServer `wire:"-"`
	GroupMemberRepo    *repo.GroupMember
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

func (h *Handler) SetServ(serv longnet.IServer) {
	h.serv = serv
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

func Message(cmd string, body any) []byte {
	msg := map[string]any{
		"event":   cmd,
		"payload": body,
	}

	data, _ := json.Marshal(msg)
	return data
}
