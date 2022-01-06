// +build wireinject

package main

import (
	"context"
	"go-chat/internal/dao"
	note2 "go-chat/internal/dao/note"
	"go-chat/internal/http/internal/handler/api/v1/article"
	"go-chat/internal/pkg/client"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/provider"
	"go-chat/internal/service/note"

	"github.com/google/wire"
	"go-chat/internal/cache"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/handler/api/v1"
	"go-chat/internal/http/internal/handler/open"
	"go-chat/internal/http/internal/router"
	"go-chat/internal/service"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewConfig,
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
	wire.Struct(new(cache.SmsCodeCache), "*"),

	// dao 数据层
	dao.NewBaseDao,
	dao.NewUsersFriendsDao,
	dao.NewGroupMemberDao,
	dao.NewUserDao,
	dao.NewGroupDao,
	wire.Struct(new(dao.TalkRecordsDao), "*"),
	wire.Struct(new(dao.TalkRecordsCodeDao), "*"),
	wire.Struct(new(dao.TalkRecordsLoginDao), "*"),
	wire.Struct(new(dao.TalkRecordsFileDao), "*"),
	wire.Struct(new(dao.GroupNoticeDao), "*"),
	dao.NewTalkSessionDao,
	dao.NewEmoticonDao,
	dao.NewTalkRecordsVoteDao,
	dao.NewFileSplitUploadDao,
	note2.NewArticleClassDao,
	note2.NewArticleAnnexDao,

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
	service.NewTalkSessionService,
	service.NewTalkMessageForwardService,
	service.NewEmoticonService,
	service.NewTalkRecordsService,
	service.NewContactService,
	service.NewContactsApplyService,
	service.NewSplitUploadService,
	service.NewIpAddressService,
	note.NewArticleService,
	note.NewArticleTagService,
	note.NewArticleClassService,
	note.NewArticleAnnexService,

	// handler 处理
	v1.NewAuthHandler,
	v1.NewCommonHandler,
	v1.NewUserHandler,
	v1.NewContactHandler,
	v1.NewContactsApplyHandler,
	v1.NewGroupHandler,
	v1.NewGroupNoticeHandler,
	v1.NewTalkHandler,
	v1.NewTalkMessageHandler,
	v1.NewUploadHandler,
	v1.NewEmoticonHandler,
	v1.NewTalkRecordsHandler,
	open.NewIndexHandler,
	article.NewAnnexHandler,
	article.NewArticleHandler,
	article.NewClassHandler,
	article.NewTagHandler,

	wire.Struct(new(handler.Handler), "*"),
	wire.Struct(new(Providers), "*"),
)

func Initialize(ctx context.Context) *Providers {
	panic(wire.Build(providerSet))
}
