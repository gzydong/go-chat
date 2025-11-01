package admin

import (
	"github.com/google/wire"
	"go-chat/internal/apis/handler/admin/system"
	"go-chat/internal/apis/handler/admin/user"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(Auth), "*"),
	wire.Struct(new(Totp), "*"),
	wire.Struct(new(system.Admin), "*"),
	wire.Struct(new(system.Role), "*"),
	wire.Struct(new(system.Resource), "*"),
	wire.Struct(new(system.Menu), "*"),
	wire.Struct(new(user.User), "*"),
)
