package mission

import (
	"github.com/google/wire"
	"go-chat/internal/mission/cron"
	"go-chat/internal/mission/queue"
	"go-chat/internal/mission/temp"
)

var CronProviderSet = wire.NewSet(
	wire.Struct(new(CronProvider), "*"),
	cron.ProviderSet,
)

var QueueProviderSet = wire.NewSet(
	wire.Struct(new(QueueProvider), "*"),
	queue.ProviderSet,
)

var MigrateProviderSet = wire.NewSet(
	wire.Struct(new(MigrateProvider), "*"),
)

var TempProviderSet = wire.NewSet(
	wire.Struct(new(temp.TestCommand), "*"),
	wire.Struct(new(TempProvider), "*"),
)
