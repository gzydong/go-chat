//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"go-chat/config"
	"go-chat/internal/provider"
	cache2 "go-chat/internal/repository/cache"
	dao2 "go-chat/internal/repository/dao"
	"go-chat/internal/service"
	handle2 "go-chat/internal/websocket/internal/handler"
	"go-chat/internal/websocket/internal/process"
	handle "go-chat/internal/websocket/internal/process/handle"
	"go-chat/internal/websocket/internal/router"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewWebsocketServer,
	router.NewRouter,

	// process
	process.NewCoroutine,
	process.NewHealth,
	process.NewWsSubscribe,
	handle.NewSubscribeConsume,

	// 缓存
	cache2.NewSession,
	cache2.NewSid,
	cache2.NewRedisLock,
	cache2.NewWsClientSession,
	cache2.NewRoom,
	cache2.NewTalkVote,
	cache2.NewRelation,

	// dao 数据层
	dao2.NewBaseDao,
	dao2.NewTalkRecordsDao,
	dao2.NewTalkRecordsVoteDao,
	dao2.NewGroupMemberDao,
	dao2.NewContactDao,

	// 服务
	service.NewBaseService,
	service.NewTalkRecordsService,
	service.NewClientService,
	service.NewGroupMemberService,
	service.NewContactService,

	// handle
	handle2.NewDefaultWebSocket,
	handle2.NewExampleWebsocket,

	wire.Struct(new(handle2.Handler), "*"),
	wire.Struct(new(Provider), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *Provider {
	panic(wire.Build(providerSet))
}
