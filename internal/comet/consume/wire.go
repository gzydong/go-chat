package consume

import (
	"github.com/google/wire"
	"go-chat/internal/comet/consume/chat"
	"go-chat/internal/comet/consume/example"
)

var ProviderSet = wire.NewSet(
	NewChatSubscribe,
	wire.Struct(new(chat.Handler), "*"),

	NewExampleSubscribe,
	example.NewHandler,
)
