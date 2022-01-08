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

	// dao
	dao.NewBaseDao,
	dao.NewFileSplitUploadDao,

	// 注册命令行
	crontab.NewCrontabCommand,
	queue.NewQueueCommand,
	other.NewOtherCommand,

	// 子命令
	crontab.NewClearTmpFileCommand,

	wire.Struct(new(cmd.Commands), "*"),
	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
