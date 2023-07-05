//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/cmd/internal/command"
	"go-chat/internal/cmd/internal/command/cron"
	"go-chat/internal/cmd/internal/command/other"
	"go-chat/internal/cmd/internal/command/queue"
	cron2 "go-chat/internal/cmd/internal/handle/cron"
	other2 "go-chat/internal/cmd/internal/handle/other"
	queue2 "go-chat/internal/cmd/internal/handle/queue"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	provider.NewEmailClient,
	provider.NewRequestClient,

	filesystem.NewFilesystem,

	// cache
	cache.NewSidStorage,
	cache.NewSequence,

	repo.NewSource,
	repo.NewSequence,

	// Crontab 命令行
	cron.NewCrontabCommand,
	cron2.NewClearTmpFile,
	cron2.NewClearArticle,
	cron2.NewClearWsCache,
	cron2.NewClearExpireServer,
	wire.Struct(new(cron.Subcommands), "*"),

	// Queue Command
	queue.NewQueueCommand,
	wire.Struct(new(queue.Subcommands), "*"),
	queue2.NewEmailHandle,

	// Other Command
	other.NewOtherCommand,
	other.NewExampleCommand,
	other.NewMigrateCommand,
	wire.Struct(new(other.Subcommands), "*"),
	other2.NewExampleHandle,

	// 服务
	wire.Struct(new(command.Commands), "*"),
	wire.Struct(new(AppProvider), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *AppProvider {
	panic(wire.Build(providerSet))
}
