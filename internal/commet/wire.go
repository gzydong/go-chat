package commet

import (
	"github.com/google/wire"
	"go-chat/internal/commet/consume"
	"go-chat/internal/commet/event"
	"go-chat/internal/commet/handler"
	"go-chat/internal/commet/process"
	"go-chat/internal/commet/router"
)

var ProviderSet = wire.NewSet(
	router.NewRouter,
	wire.Struct(new(handler.Handler), "*"),

	// process
	wire.Struct(new(process.SubServers), "*"),
	process.NewServer,
	process.NewHealthSubscribe,
	process.NewMessageSubscribe,

	handler.ProviderSet,
	event.ProviderSet,
	consume.ProviderSet,

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)
