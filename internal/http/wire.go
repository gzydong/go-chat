//go:build wireinject
// +build wireinject

package main

import (
	"go-chat/config"
	"go-chat/internal/http/internal/handler/admin"
	"go-chat/internal/http/internal/handler/open"
	"go-chat/internal/http/internal/handler/web"
	"go-chat/internal/logic"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	note3 "go-chat/internal/repository/repo/note"
	organize3 "go-chat/internal/repository/repo/organize"
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
	provider.NewFilesystem,
	provider.NewRequestClient,

	// 注册路由
	router.NewRouter,
	wire.Struct(new(web.Handler), "*"),
	wire.Struct(new(admin.Handler), "*"),
	wire.Struct(new(open.Handler), "*"),
	wire.Struct(new(handler.Handler), "*"),

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)

var daoProviderSet = wire.NewSet(
	repo.NewContact,
	repo.NewContactGroup,
	repo.NewGroupMember,
	repo.NewUsers,
	repo.NewGroup,
	repo.NewGroupApply,
	repo.NewTalkRecords,
	repo.NewGroupNotice,
	repo.NewTalkSession,
	repo.NewEmoticon,
	repo.NewTalkRecordsVote,
	repo.NewFileSplitUpload,
	note3.NewArticleClass,
	note3.NewArticleAnnex,
	organize3.NewDepartment,
	organize3.NewOrganize,
	organize3.NewPosition,
	repo.NewRobot,
	repo.NewTest,
	repo.NewSequence,
)

var serviceProviderSet = wire.NewSet(
	service.NewBaseService,
	service.NewUserService,
	service.NewSmsService,
	service.NewTalkService,
	service.NewTalkMessageService,
	service.NewGroupService,
	service.NewGroupMemberService,
	service.NewGroupNoticeService,
	service.NewGroupApplyService,
	service.NewTalkSessionService,
	service.NewEmoticonService,
	service.NewTalkRecordsService,
	service.NewContactService,
	service.NewContactApplyService,
	service.NewContactGroupService,
	service.NewSplitUploadService,
	service.NewIpAddressService,
	service.NewAuthPermissionService,
	service.NewMessageService,
	note.NewArticleService,
	note.NewArticleTagService,
	note.NewArticleClassService,
	note.NewArticleAnnexService,
	organize.NewOrganizeDeptService,
	organize.NewOrganizeService,
	organize.NewPositionService,
	service.NewTemplateService,
	service.NewTalkAuthService,
	logic.NewMessageForwardLogic,
)

func Initialize(conf *config.Config) *AppProvider {
	panic(
		wire.Build(
			providerSet,
			cache.ProviderSet,  // 注入 Cache 依赖
			daoProviderSet,     // 注入 Dao 依赖
			serviceProviderSet, // 注入 Service 依赖
			web.ProviderSet,    // 注入 Web Handler 依赖
			admin.ProviderSet,  // 注入 Admin Handler 依赖
			open.ProviderSet,   // 注入 Open Handler 依赖
		),
	)
}
