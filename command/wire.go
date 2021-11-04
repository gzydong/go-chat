// +build wireinject

package main

import (
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/provider"
	"gorm.io/gorm"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewLogger,
	provider.NewRedisClient,
	provider.NewMySQLClient,
	provider.NewHttp,
)

func Initialize(conf *config.Config) *gorm.DB {
	panic(wire.Build(providerSet))
}
