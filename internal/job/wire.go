//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/job/internal/command"
	"go-chat/internal/job/internal/command/cron"
	"go-chat/internal/job/internal/command/other"
	"go-chat/internal/job/internal/command/queue"
	crontab "go-chat/internal/job/internal/handle/crontab"
	other2 "go-chat/internal/job/internal/handle/other"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewConfig,
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
	crontab.NewClearTmpFile,
	crontab.NewClearArticle,
	crontab.NewClearWsCacheHandle,
	crontab.NewClearExpireServer,
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
	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
