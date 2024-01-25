package event

import (
	"github.com/google/wire"
	"go-chat/internal/comet/handler/event/chat"
	"go-chat/internal/comet/handler/event/example"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(ChatEvent), "*"),

	wire.Struct(new(chat.Handler), "*"),

	wire.Struct(new(ExampleEvent), "*"),
	example.NewHandler,
)
