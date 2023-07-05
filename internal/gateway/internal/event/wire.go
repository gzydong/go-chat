package event

import (
	"github.com/google/wire"
	"go-chat/internal/gateway/internal/event/chat"
	"go-chat/internal/gateway/internal/event/example"
)

var ProviderSet = wire.NewSet(
	NewChatEvent,
	chat.NewHandler,

	NewExampleEvent,
	example.NewHandler,
)
