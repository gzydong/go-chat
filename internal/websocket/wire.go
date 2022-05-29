//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/dao"
	"go-chat/internal/provider"
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
	cache.NewSession,
	cache.NewSid,
	cache.NewRedisLock,
	cache.NewWsClientSession,
	cache.NewRoom,
	cache.NewTalkVote,
	cache.NewRelation,

	// dao 数据层
	dao.NewBaseDao,
	dao.NewTalkRecordsDao,
	dao.NewTalkRecordsVoteDao,
	dao.NewGroupMemberDao,
	dao.NewContactDao,

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
