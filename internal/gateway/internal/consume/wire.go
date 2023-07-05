package consume

import (
	"github.com/google/wire"
	"go-chat/internal/gateway/internal/consume/chat"
	"go-chat/internal/gateway/internal/consume/example"
)

var ProviderSet = wire.NewSet(
	NewChatSubscribe,
	chat.NewHandler,

	NewExampleSubscribe,
	example.NewHandler,
)
