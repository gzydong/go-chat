package handler

import (
	"github.com/google/wire"
	"go-chat/internal/api/handler/admin"
	"go-chat/internal/api/handler/open"
	"go-chat/internal/api/handler/web"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(web.Handler), "*"),
	wire.Struct(new(admin.Handler), "*"),
	wire.Struct(new(open.Handler), "*"),
	wire.Struct(new(Handler), "*"),
)
