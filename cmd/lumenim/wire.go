//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/business"
	"go-chat/internal/comet"
	"go-chat/internal/mission"
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
	provider.NewBase64Captcha,
	provider.NewIpAddressClient,
	wire.Struct(new(provider.Providers), "*"),

	cache.ProviderSet,    // 注入 Cache 依赖
	repo.ProviderSet,     // 注入 Repo 依赖
	business.ProviderSet, // 注入 Logic 依赖
	service.ProviderSet,  // 注入 Service 依赖
)

func NewHttpInjector(conf *config.Config) *apis.AppProvider {
	panic(
		wire.Build(
			providerSet,
			apis.ProviderSet,
		),
	)
}

func NewCommetInjector(conf *config.Config) *comet.AppProvider {
	panic(
		wire.Build(
			providerSet,
			comet.ProviderSet,
		),
	)
}

func NewCronInjector(conf *config.Config) *mission.CronProvider {
	panic(
		wire.Build(
			providerSet,
			mission.CronProviderSet,
		),
	)
}

func NewQueueInjector(conf *config.Config) *mission.QueueProvider {
	panic(
		wire.Build(
			providerSet,
			mission.QueueProviderSet,
		),
	)
}

func NewOtherInjector(conf *config.Config) *mission.TempProvider {
	panic(
		wire.Build(
			providerSet,
			mission.TempProviderSet,
		),
	)
}

func NewMigrateInjector(conf *config.Config) *mission.MigrateProvider {
	panic(
		wire.Build(
			providerSet,
			mission.MigrateProviderSet,
		),
	)
}
