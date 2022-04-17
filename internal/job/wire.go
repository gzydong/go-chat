//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/job/internal/cmd"
	"go-chat/internal/job/internal/cmd/cron"
	"go-chat/internal/job/internal/cmd/other"
	"go-chat/internal/job/internal/cmd/queue"
	crontab "go-chat/internal/job/internal/handle/crontab"
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

	// 命令行
	cron.NewCrontabCommand,

	// Queue Command
	queue.NewQueueCommand,
	wire.Struct(new(queue.Subcommands), "*"),

	// Other Command
	other.NewOtherCommand,
	other.NewTestCommand,
	wire.Struct(new(other.Subcommands), "*"),

	// Handle
	crontab.NewClearTmpFile,
	crontab.NewClearArticle,
	crontab.NewClearWsCacheHandle,
	crontab.NewClearExpireServer,
	wire.Struct(new(cron.Handles), "*"),

	// 服务
	wire.Struct(new(cmd.Commands), "*"),
	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
