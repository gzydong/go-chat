package handler

import (
	"github.com/google/wire"
	"go-chat/internal/httpapi/handler/admin"
	"go-chat/internal/httpapi/handler/open"
	"go-chat/internal/httpapi/handler/web"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(web.Handler), "*"),
	wire.Struct(new(admin.Handler), "*"),
	wire.Struct(new(open.Handler), "*"),
	wire.Struct(new(Handler), "*"),
)
