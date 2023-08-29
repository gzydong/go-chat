//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/handler/admin"
	"go-chat/internal/http/internal/handler/open"
	"go-chat/internal/http/internal/handler/web"
	"go-chat/internal/http/internal/router"
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

	// 注册路由
	router.NewRouter,

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)

func Initialize(conf *config.Config) *AppProvider {
	panic(
		wire.Build(
			providerSet,
			cache.ProviderSet,   // 注入 Cache 依赖
			repo.ProviderSet,    // 注入 Repo 依赖
			logic.ProviderSet,   // 注入 Logic 依赖
			service.ProviderSet, // 注入 Service 依赖
			handler.ProviderSet, // 注入 Handler
			web.ProviderSet,     // 注入 Web Handler 依赖
			admin.ProviderSet,   // 注入 Admin Handler 依赖
			open.ProviderSet,    // 注入 Open Handler 依赖
		),
	)
}
