package consume

import (
	"github.com/google/wire"
	"go-chat/internal/commet/consume/chat"
	"go-chat/internal/commet/consume/example"
)

var ProviderSet = wire.NewSet(
	NewChatSubscribe,
	wire.Struct(new(chat.Handler), "*"),

	NewExampleSubscribe,
	example.NewHandler,
)
