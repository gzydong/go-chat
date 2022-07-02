package admin

import (
	"github.com/google/wire"
	v1 "go-chat/internal/http/internal/handler/admin/v1"
)

var ProviderSet = wire.NewSet(
	v1.NewIndex,
)
