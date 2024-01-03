package mission

import (
	"github.com/google/wire"
	"go-chat/internal/mission/cron"
	"go-chat/internal/mission/queue"
	"go-chat/internal/mission/temp"
)

var CronProviderSet = wire.NewSet(
	wire.Struct(new(CronProvider), "*"),
	wire.Struct(new(Crontab), "*"),
	wire.Struct(new(cron.ClearArticle), "*"),
	wire.Struct(new(cron.ClearTmpFile), "*"),
	cron.NewClearWsCache,
	cron.NewClearExpireServer,
)

var QueueProviderSet = wire.NewSet(
	wire.Struct(new(QueueProvider), "*"),
	wire.Struct(new(QueueJobs), "*"),
	wire.Struct(new(queue.ExampleQueue), "*"),
)

var MigrateProviderSet = wire.NewSet(
	wire.Struct(new(MigrateProvider), "*"),
)

var TempProviderSet = wire.NewSet(
	wire.Struct(new(temp.TestCommand), "*"),
	wire.Struct(new(TempProvider), "*"),
)
