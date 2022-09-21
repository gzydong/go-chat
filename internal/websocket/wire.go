//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"go-chat/config"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/dao"
	"go-chat/internal/service"
	"go-chat/internal/websocket/internal/event"
	"go-chat/internal/websocket/internal/event/chat"
	"go-chat/internal/websocket/internal/handler"
	"go-chat/internal/websocket/internal/process"
	"go-chat/internal/websocket/internal/process/consume"
	"go-chat/internal/websocket/internal/process/server"
	"go-chat/internal/websocket/internal/router"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewWebsocketServer,

	// 路由
	router.NewRouter,

	// process
	wire.Struct(new(process.SubServers), "*"),
	process.NewServer,
	server.NewHealthSubscribe,
	server.NewMessageSubscribe,
	consume.NewDefaultSubscribe,
	consume.NewExampleSubscribe,

	// 缓存
	cache.NewSessionStorage,
	cache.NewSid,
	cache.NewRedisLock,
	cache.NewClientStorage,
	cache.NewRoomStorage,
	cache.NewTalkVote,
	cache.NewRelation,
	cache.NewContactRemark,

	// dao 数据层
	dao.NewBaseDao,
	dao.NewTalkRecordsDao,
	dao.NewTalkRecordsVoteDao,
	dao.NewGroupMemberDao,
	dao.NewContactDao,

	chat.NewHandler,

	event.NewDefaultEvent,
	event.NewExampleEvent,

	// 服务
	service.NewBaseService,
	service.NewTalkRecordsService,
	service.NewClientService,
	service.NewGroupMemberService,
	service.NewContactService,

	// handle
	handler.NewDefaultChannel,
	handler.NewExampleChannel,

	wire.Struct(new(handler.Handler), "*"),
	wire.Struct(new(AppProvider), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *AppProvider {
	panic(wire.Build(providerSet))
}
