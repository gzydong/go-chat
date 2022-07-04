package admin

import (
	"github.com/google/wire"
	"go-chat/internal/http/internal/handler/admin/v1"
)

var ProviderSet = wire.NewSet(
	v1.NewIndex,
	v1.NewAuth,

	wire.Struct(new(V1), "*"),
	wire.Struct(new(V2), "*"),
)
