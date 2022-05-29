//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"go-chat/config"
	"go-chat/internal/dao"
	note2 "go-chat/internal/dao/note"
	organize2 "go-chat/internal/dao/organize"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
	"go-chat/internal/service/note"
	"go-chat/internal/service/organize"

	"github.com/google/wire"
	"go-chat/internal/cache"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/router"
	"go-chat/internal/service"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	provider.NewHttpServer,
	client.NewHttpClient,

	// 注册路由
	router.NewRouter,

	// other
	filesystem.NewFilesystem,

	// 缓存
	cache.NewSession,
	cache.NewSid,
	cache.NewUnreadTalkCache,
	cache.NewRedisLock,
	cache.NewWsClientSession,
	cache.NewLastMessage,
	cache.NewTalkVote,
	cache.NewRoom,
	cache.NewRelation,
	cache.NewSmsCodeCache,

	// dao 数据层
	dao.NewBaseDao,
	dao.NewContactDao,
	dao.NewGroupMemberDao,
	dao.NewUserDao,
	dao.NewGroupDao,
	dao.NewGroupApply,
	dao.NewTalkRecordsDao,
	dao.NewGroupNoticeDao,
	dao.NewTalkSessionDao,
	dao.NewEmoticonDao,
	dao.NewTalkRecordsVoteDao,
	dao.NewFileSplitUploadDao,
	note2.NewArticleClassDao,
	note2.NewArticleAnnexDao,
	organize2.NewDepartmentDao,
	organize2.NewOrganizeDao,
	organize2.NewPositionDao,

	// 服务
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

	// Handler
	wire.Struct(new(handler.ApiHandler), "*"),
	wire.Struct(new(handler.AdminHandler), "*"),
	wire.Struct(new(handler.OpenHandler), "*"),
	wire.Struct(new(handler.Handler), "*"),

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)

func Initialize(ctx context.Context, conf *config.Config) *AppProvider {
	panic(wire.Build(providerSet, handler.ProviderSet))
}
