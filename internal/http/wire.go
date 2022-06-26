//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"go-chat/config"
	"go-chat/internal/pkg/client"
	"go-chat/internal/provider"
	cache2 "go-chat/internal/repository/cache"
	dao2 "go-chat/internal/repository/dao"
	note3 "go-chat/internal/repository/dao/note"
	organize3 "go-chat/internal/repository/dao/organize"
	"go-chat/internal/service/note"
	"go-chat/internal/service/organize"

	"github.com/google/wire"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/router"
	"go-chat/internal/service"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	provider.NewEmailClient,
	provider.NewHttpServer,
	provider.NewFilesystem,
	client.NewHttpClient,

	// 注册路由
	router.NewRouter,
	wire.Struct(new(handler.ApiHandler), "*"),
	wire.Struct(new(handler.AdminHandler), "*"),
	wire.Struct(new(handler.OpenHandler), "*"),
	wire.Struct(new(handler.Handler), "*"),

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)

var cacheProviderSet = wire.NewSet(
	cache2.NewSession,
	cache2.NewSid,
	cache2.NewUnreadTalkCache,
	cache2.NewRedisLock,
	cache2.NewWsClientSession,
	cache2.NewLastMessage,
	cache2.NewTalkVote,
	cache2.NewRoom,
	cache2.NewRelation,
	cache2.NewSmsCodeCache,
)

var daoProviderSet = wire.NewSet(
	dao2.NewBaseDao,
	dao2.NewContactDao,
	dao2.NewGroupMemberDao,
	dao2.NewUserDao,
	dao2.NewGroupDao,
	dao2.NewGroupApply,
	dao2.NewTalkRecordsDao,
	dao2.NewGroupNoticeDao,
	dao2.NewTalkSessionDao,
	dao2.NewEmoticonDao,
	dao2.NewTalkRecordsVoteDao,
	dao2.NewFileSplitUploadDao,
	note3.NewArticleClassDao,
	note3.NewArticleAnnexDao,
	organize3.NewDepartmentDao,
	organize3.NewOrganizeDao,
	organize3.NewPositionDao,
)

var serviceProviderSet = wire.NewSet(
	service.NewBaseService,
	service.NewUserService,
	service.NewSmsService,
	service.NewTalkService,
	service.NewTalkMessageService,
	service.NewClientService,
	service.NewGroupService,
	service.NewGroupMemberService,
	service.NewGroupNoticeService,
	service.NewGroupApplyService,
	service.NewTalkSessionService,
	service.NewTalkMessageForwardService,
	service.NewEmoticonService,
	service.NewTalkRecordsService,
	service.NewContactService,
	service.NewContactsApplyService,
	service.NewSplitUploadService,
	service.NewIpAddressService,
	service.NewAuthPermissionService,
	note.NewArticleService,
	note.NewArticleTagService,
	note.NewArticleClassService,
	note.NewArticleAnnexService,
	organize.NewOrganizeDeptService,
	organize.NewOrganizeService,
	organize.NewPositionService,
	service.NewTemplateService,
)

func Initialize(ctx context.Context, conf *config.Config) *AppProvider {
	panic(
		wire.Build(
			providerSet,
			cacheProviderSet,   // 注入 Cache 依赖
			daoProviderSet,     // 注入 Dao 依赖
			serviceProviderSet, // 注入 Service 依赖
			handler.ProviderSet,
		),
	)
}
