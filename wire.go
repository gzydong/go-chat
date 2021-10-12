// +build wireinject

package main

import (
	"context"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/repository"

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
	connect.MysqlConnect,
	connect.NewHttp,
	router.NewRouter,

	filesystem.NewFilesystem,

	// 缓存
	cache.NewServerRun,
	wire.Struct(new(cache.WsClient), "*"),
	wire.Struct(new(cache.AuthTokenCache), "*"),
	wire.Struct(new(cache.SmsCodeCache), "*"),
	wire.Struct(new(cache.RedisLock), "*"),

	// handler 处理
	wire.Struct(new(v1.Common), "*"),
	wire.Struct(new(v1.Auth), "*"),
	wire.Struct(new(v1.User), "*"),
	wire.Struct(new(v1.Upload), "*"),
	wire.Struct(new(v1.Download), "*"),
	wire.Struct(new(v1.TalkMessage), "*"),
	wire.Struct(new(open.Index), "*"),
	wire.Struct(new(ws.WebSocket), "*"),
	wire.Struct(new(handler.Handler), "*"),

	// repo
	wire.Struct(new(repository.UserRepository), "*"),
	wire.Struct(new(repository.TalkRecordsRepo), "*"),
	wire.Struct(new(repository.TalkRecordsCodeRepo), "*"),
	wire.Struct(new(repository.TalkRecordsLoginRepo), "*"),
	wire.Struct(new(repository.TalkRecordsFileRepo), "*"),
	wire.Struct(new(repository.TalkRecordsVoteRepo), "*"),

	// 服务
	wire.Struct(new(service.ClientService), "*"),
	wire.Struct(new(service.UserService), "*"),
	wire.Struct(new(service.SocketService), "*"),
	wire.Struct(new(service.SmsService), "*"),
	wire.Struct(new(service.TalkMessageService), "*"),
	wire.Struct(new(Service), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *Service {
	panic(wire.Build(providerSet))
}
