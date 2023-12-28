//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/commet"
	"go-chat/internal/job"
	"go-chat/internal/logic"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	provider.NewEmailClient,
	provider.NewFilesystem,
	provider.NewRequestClient,
	provider.NewBase64Captcha,
	wire.Struct(new(provider.Providers), "*"),

	cache.ProviderSet,   // 注入 Cache 依赖
	repo.ProviderSet,    // 注入 Repo 依赖
	logic.ProviderSet,   // 注入 Logic 依赖
	service.ProviderSet, // 注入 Service 依赖
)

func NewHttpInjector(conf *config.Config) *apis.AppProvider {
	panic(
		wire.Build(
			providerSet,
			apis.ProviderSet,
		),
	)
}

func NewCommetInjector(conf *config.Config) *commet.AppProvider {
	panic(
		wire.Build(
			providerSet,
			commet.ProviderSet,
		),
	)
}

func NewCronInjector(conf *config.Config) *job.CronProvider {
	panic(
		wire.Build(
			providerSet,
			job.CronProviderSet,
		),
	)
}

func NewQueueInjector(conf *config.Config) *job.QueueProvider {
	panic(
		wire.Build(
			providerSet,
			job.QueueProviderSet,
		),
	)
}

func NewOtherInjector(conf *config.Config) *job.TempProvider {
	panic(
		wire.Build(
			providerSet,
			job.TempProviderSet,
		),
	)
}

func NewMigrateInjector(conf *config.Config) *job.MigrateProvider {
	panic(
		wire.Build(
			providerSet,
			job.MigrateProviderSet,
		),
	)
}
