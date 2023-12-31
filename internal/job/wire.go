package job

import (
	"github.com/google/wire"
	"go-chat/internal/job/cron"
	"go-chat/internal/job/queue"
	"go-chat/internal/job/temp"
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
	wire.Struct(new(Queue), "*"),
	wire.Struct(new(queue.ExampleQueue), "*"),
)

var MigrateProviderSet = wire.NewSet(
	wire.Struct(new(MigrateProvider), "*"),
)

var TempProviderSet = wire.NewSet(
	wire.Struct(new(temp.TestCommand), "*"),
	wire.Struct(new(TempProvider), "*"),
)
