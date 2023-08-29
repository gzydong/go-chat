package event

import (
	"github.com/google/wire"
	"go-chat/internal/gateway/internal/event/chat"
	"go-chat/internal/gateway/internal/event/example"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(ChatEvent), "*"),
	chat.NewHandler,

	wire.Struct(new(ExampleEvent), "*"),
	example.NewHandler,
)
