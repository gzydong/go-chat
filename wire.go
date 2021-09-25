// +build wireinject

package main

import (
	"context"
	"net/http"

	"github.com/google/wire"
	"go-chat/app/cache"
	"go-chat/app/http/handler"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/router"
	"go-chat/app/service"
	"go-chat/config"
	"go-chat/connect"
)

var providerSet = wire.NewSet(
	// 连接信息
	connect.RedisConnect,
	connect.NewHttp,
	router.NewRouter,

	// 缓存
	cache.NewServerRun,
	wire.Struct(new(cache.WsClient), "*"),

	// handler 处理
	wire.Struct(new(v1.Auth), "*"),
	wire.Struct(new(v1.User), "*"),
	wire.Struct(new(v1.Download), "*"),
	wire.Struct(new(open.Index), "*"),
	wire.Struct(new(ws.Ws), "*"),
	wire.Struct(new(handler.Handler), "*"),

	// 服务
	wire.Struct(new(service.ClientService), "*"),
	wire.Struct(new(service.UserService), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *http.Server {
	panic(wire.Build(providerSet))
}
