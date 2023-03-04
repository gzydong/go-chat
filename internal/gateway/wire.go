//go:build wireinject
// +build wireinject

package main

import (
	"go-chat/config"
	consume2 "go-chat/internal/gateway/internal/consume"
	"go-chat/internal/gateway/internal/event"
	"go-chat/internal/gateway/internal/event/chat"
	"go-chat/internal/gateway/internal/handler"
	"go-chat/internal/gateway/internal/process"
	"go-chat/internal/gateway/internal/router"
	"go-chat/internal/logic"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewFilesystem,
	provider.NewEmailClient,
	provider.NewProviders,

	// 路由
	router.NewRouter,

	// process
	wire.Struct(new(process.SubServers), "*"),
	process.NewServer,
	process.NewHealthSubscribe,
	process.NewMessageSubscribe,
	consume2.NewChatSubscribe,
	consume2.NewExampleSubscribe,

	// 缓存
	cache.NewTokenSessionStorage,
	cache.NewSidStorage,
	cache.NewRedisLock,
	cache.NewClientStorage,
	cache.NewRoomStorage,
	cache.NewTalkVote,
	cache.NewRelation,
	cache.NewContactRemark,
	cache.NewSequence,
	cache.NewUnreadStorage,
	cache.NewMessageStorage,

	// dao 数据层
	repo.NewTalkRecords,
	repo.NewTalkRecordsVote,
	repo.NewGroupMember,
	repo.NewContact,
	repo.NewFileSplitUpload,
	repo.NewSequence,

	logic.NewMessageForwardLogic,

	chat.NewHandler,

	event.NewChatEvent,
	event.NewExampleEvent,

	// 服务
	service.NewBaseService,
	service.NewTalkRecordsService,
	service.NewGroupMemberService,
	service.NewContactService,
	service.NewMessageService,

	// handle
	handler.NewChatChannel,
	handler.NewExampleChannel,

	wire.Struct(new(handler.Handler), "*"),
	wire.Struct(new(AppProvider), "*"),
)

func Initialize(conf *config.Config) *AppProvider {
	panic(wire.Build(providerSet))
}
