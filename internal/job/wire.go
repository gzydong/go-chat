// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/pkg/client"
	"go-chat/internal/provider"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewConfig,
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	client.NewHttpClient,
)

func Initialize(ctx context.Context) *config.Config {
	panic(wire.Build(providerSet))
}
