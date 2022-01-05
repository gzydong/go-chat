// +build wireinject

package main

import (
	"context"
	"go-chat/app/cache"
	"go-chat/app/dao"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/service"
	handle2 "go-chat/app/websocket/internal/handler"
	"go-chat/app/websocket/internal/process"
	handle "go-chat/app/websocket/internal/process/handle"
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

	// process
	process.NewProcess,
	process.NewClearGarbage,
	process.NewImHeartbeat,
	process.NewServer,
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
	dao.NewUsersFriendsDao,
	filesystem.NewFilesystem,

	// 服务
	service.NewBaseService,
	service.NewTalkRecordsService,
	service.NewClientService,
	service.NewGroupMemberService,
	service.NewContactService,

	// handle
	handle2.NewDefaultWebSocket,

	wire.Struct(new(handle2.Handler), "*"),
	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
