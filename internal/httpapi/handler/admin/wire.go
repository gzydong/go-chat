package admin

import (
	"github.com/google/wire"
	v12 "go-chat/internal/httpapi/handler/admin/v1"
)

var ProviderSet = wire.NewSet(
	v12.NewIndex,
	v12.NewAuth,

	wire.Struct(new(V1), "*"),
	wire.Struct(new(V2), "*"),
)
