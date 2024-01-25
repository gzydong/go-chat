package comet

import (
	"github.com/google/wire"
	"go-chat/internal/comet/consume"
	"go-chat/internal/comet/handler"
	"go-chat/internal/comet/handler/event"
	"go-chat/internal/comet/process"
	"go-chat/internal/comet/process/queue"
	"go-chat/internal/comet/router"
	"go-chat/internal/pkg/core/socket"
)

var ProviderSet = wire.NewSet(
	router.NewRouter,
	socket.NewRoomStorage,

	wire.Struct(new(handler.Handler), "*"),

	// process
	wire.Struct(new(process.SubServers), "*"),
	process.NewServer,
	process.NewHealthSubscribe,
	process.NewMessageSubscribe,
	wire.Struct(new(process.QueueSubscribe), "*"),
	wire.Struct(new(queue.GlobalMessage), "*"),
	wire.Struct(new(queue.LocalMessage), "*"),
	wire.Struct(new(queue.RoomControl), "*"),

	handler.ProviderSet,
	event.ProviderSet,
	consume.ProviderSet,

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)
