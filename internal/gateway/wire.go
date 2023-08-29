//go:build wireinject
// +build wireinject

package main

import (
	"go-chat/config"
	"go-chat/internal/gateway/internal/consume"
	"go-chat/internal/gateway/internal/event"
	"go-chat/internal/gateway/internal/handler"
	"go-chat/internal/gateway/internal/process"
	"go-chat/internal/gateway/internal/router"
	"go-chat/internal/logic"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/organize"
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

	// 数据层
	repo.NewSource,
	repo.NewTalkRecords,
	repo.NewTalkRecordsVote,
	repo.NewGroupMember,
	repo.NewContact,
	repo.NewFileSplitUpload,
	repo.NewSequence,
	organize.NewOrganize,
	repo.NewRobot,

	logic.NewMessageForwardLogic,

	// 服务
	service.NewTalkRecordsService,
	service.NewGroupMemberService,
	service.NewContactService,

	wire.Struct(new(service.MessageService), "*"),
	wire.Bind(new(service.IMessageService), new(*service.MessageService)),

	wire.Struct(new(handler.Handler), "*"),
	wire.Struct(new(AppProvider), "*"),
)

func Initialize(conf *config.Config) *AppProvider {
	panic(wire.Build(providerSet, cache.ProviderSet, handler.ProviderSet, event.ProviderSet, consume.ProviderSet))
}
