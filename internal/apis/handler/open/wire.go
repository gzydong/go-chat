package open

import (
	"github.com/google/wire"
	"go-chat/internal/apis/handler/open/v1"
)

var ProviderSet = wire.NewSet(
	v1.NewIndex,

	wire.Struct(new(V1), "*"),
)
