//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/comet"
	"go-chat/internal/logic"
	"go-chat/internal/mission"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.ProviderSet,
	cache.ProviderSet,   // 注入 Cache 依赖
	repo.ProviderSet,    // 注入 Repo 依赖
	logic.ProviderSet,   // 注入 Logic 依赖
	service.ProviderSet, // 注入 Service 依赖
)

func NewHttpInjector(c *config.Config) *apis.Provider {
	panic(
		wire.Build(
			providerSet,
			apis.ProviderSet,
		),
	)
}

func NewCometInjector(c *config.Config) *comet.Provider {
	panic(
		wire.Build(
			providerSet,
			comet.ProviderSet,
		),
	)
}

func NewCronInjector(c *config.Config) *mission.CronProvider {
	panic(
		wire.Build(
			providerSet,
			mission.CronProviderSet,
		),
	)
}

func NewQueueInjector(c *config.Config) *mission.QueueProvider {
	panic(
		wire.Build(
			providerSet,
			mission.QueueProviderSet,
		),
	)
}

func NewTempInjector(c *config.Config) *mission.TempProvider {
	panic(
		wire.Build(
			providerSet,
			mission.TempProviderSet,
		),
	)
}

func NewMigrateInjector(c *config.Config) *mission.MigrateProvider {
	panic(
		wire.Build(
			providerSet,
			mission.MigrateProviderSet,
		),
	)
}
