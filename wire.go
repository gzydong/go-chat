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

	// other
	filesystem.NewFilesystem,

	// 缓存
	cache.NewServerRun,
	wire.Struct(new(cache.WsClient), "*"),
	wire.Struct(new(cache.AuthTokenCache), "*"),
	wire.Struct(new(cache.SmsCodeCache), "*"),
	wire.Struct(new(cache.RedisLock), "*"),

	// repo
	wire.Struct(new(repository.UserRepository), "*"),
	wire.Struct(new(repository.TalkRecordsRepo), "*"),
	wire.Struct(new(repository.TalkRecordsCodeRepo), "*"),
	wire.Struct(new(repository.TalkRecordsLoginRepo), "*"),
	wire.Struct(new(repository.TalkRecordsFileRepo), "*"),
	wire.Struct(new(repository.TalkRecordsVoteRepo), "*"),

	// 服务
	service.NewUserService,
	service.NewSmsService,
	service.NewTalkMessageService,
	service.NewClientService,
	wire.Struct(new(service.SocketService), "*"),
	wire.Struct(new(Service), "*"),

	// handler 处理
	v1.NewAuthHandler,
	v1.NewCommonHandler,
	v1.NewUserHandler,
	v1.NewTalkHandler,
	v1.NewTalkMessageHandler,
	v1.NewUploadHandler,
	v1.NewDownloadHandler,
	v1.NewEmoticonHandler,
	open.NewIndexHandler,
	ws.NewWebSocketHandler,
	wire.Struct(new(handler.Handler), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *Service {
	panic(wire.Build(providerSet))
}
