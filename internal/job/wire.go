// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/internal/dao"
	"go-chat/internal/job/internal/cmd"
	"go-chat/internal/job/internal/cmd/crontab"
	"go-chat/internal/job/internal/cmd/other"
	"go-chat/internal/job/internal/cmd/queue"
	crontab2 "go-chat/internal/job/internal/handle/crontab"
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

	// Dao
	dao.NewBaseDao,

	// 命令行
	crontab.NewCrontabCommand,
	queue.NewQueueCommand,
	other.NewOtherCommand,

	// 子命令
	crontab.NewClearTmpFileCommand,
	crontab.NewClearArticleCommand,

	// Handle
	crontab2.NewClearTmpFile,
	crontab2.NewClearArticle,

	// 服务
	wire.Struct(new(cmd.Commands), "*"),
	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
