package handler

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(ChatChannel), "*"),
	wire.Struct(new(ExampleChannel), "*"),
)
