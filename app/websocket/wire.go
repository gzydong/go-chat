// +build wireinject

package main

import (
	"context"
	"go-chat/app/websocket/internal/router"
	"go-chat/provider"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewConfig,
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewWebsocketServer,
	router.NewRouter,

	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
