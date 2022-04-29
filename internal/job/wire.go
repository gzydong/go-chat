//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/job/internal/command"
	"go-chat/internal/job/internal/command/cron"
	"go-chat/internal/job/internal/command/other"
	"go-chat/internal/job/internal/command/queue"
	cron2 "go-chat/internal/job/internal/handle/cron"
	other2 "go-chat/internal/job/internal/handle/other"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	client.NewHttpClient,

	filesystem.NewFilesystem,

	// cache
	cache.NewSid,

	// Dao
	dao.NewBaseDao,

	// Crontab 命令行
	cron.NewCrontabCommand,
	cron2.NewClearTmpFile,
	cron2.NewClearArticle,
	cron2.NewClearWsCacheHandle,
	cron2.NewClearExpireServer,
	wire.Struct(new(cron.Handles), "*"),

	// Queue Command
	queue.NewQueueCommand,
	wire.Struct(new(queue.Subcommands), "*"),

	// Other Command
	other.NewOtherCommand,
	other.NewExampleCommand,
	wire.Struct(new(other.Subcommands), "*"),
	other2.NewExampleHandle,

	// 服务
	wire.Struct(new(command.Commands), "*"),
	wire.Struct(new(Provider), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *Provider {
	panic(wire.Build(providerSet))
}
