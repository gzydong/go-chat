package event

import (
	"github.com/google/wire"
	"go-chat/internal/commet/event/chat"
	"go-chat/internal/commet/event/example"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(ChatEvent), "*"),
	chat.NewHandler,

	wire.Struct(new(ExampleEvent), "*"),
	example.NewHandler,
)
