package admin

import (
	"github.com/google/wire"
	v12 "go-chat/internal/apis/handler/admin/v1"
)

var ProviderSet = wire.NewSet(
	v12.NewIndex,
	wire.Struct(new(v12.Auth), "*"),

	wire.Struct(new(V1), "*"),
	wire.Struct(new(V2), "*"),
)
