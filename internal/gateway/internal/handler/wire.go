package handler

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewChatChannel,
	NewExampleChannel,
)
