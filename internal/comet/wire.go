package comet

import (
	"github.com/google/wire"
	"github.com/gzydong/go-chat/internal/comet/consume"
)

var ProviderSet = wire.NewSet(

	wire.Struct(new(Subscribe), "*"),
	wire.Struct(new(Handler), "*"),
	wire.Struct(new(Heartbeat), "*"),
	wire.Struct(new(consume.Handler), "*"),

	wire.Struct(new(Server), "*"),
	wire.Struct(new(Provider), "*"),
)
