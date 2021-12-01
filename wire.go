// +build wireinject

package main

import (
	"context"
	"go-chat/app/dao"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/process"
	"go-chat/app/process/handle"
	"go-chat/provider"

	"github.com/google/wire"
	"go-chat/app/cache"
	"go-chat/app/http/handler"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/router"
	"go-chat/app/service"
)

var providerSet = wire.NewSet(
	// 基础服务
	provider.NewConfig,
	provider.NewLogger,
	provider.NewMySQLClient,
	provider.NewRedisClient,
	provider.NewHttpClient,
	provider.NewHttpServer,

	// 注册路由
	router.NewRouter,

	// other
	filesystem.NewFilesystem,

	// 缓存
	cache.NewSession,
	cache.NewServerRun,
	cache.NewUnreadTalkCache,
	cache.NewRedisLock,
	cache.NewWsClientSession,
	cache.NewLastMessage,
	cache.NewTalkVote,
	cache.NewGroupRoom,
	wire.Struct(new(cache.SmsCodeCache), "*"),

	// dao 数据层
	dao.NewBaseDao,
	dao.NewUsersFriends,
	wire.Struct(new(dao.UserDao), "*"),
	wire.Struct(new(dao.TalkRecordsDao), "*"),
	wire.Struct(new(dao.TalkRecordsCodeDao), "*"),
	wire.Struct(new(dao.TalkRecordsLoginDao), "*"),
	wire.Struct(new(dao.TalkRecordsFileDao), "*"),
	wire.Struct(new(dao.GroupDao), "*"),
	wire.Struct(new(dao.GroupNoticeDao), "*"),
	dao.NewTalkListDao,
	dao.NewEmoticonDao,
	dao.NewTalkRecordsVoteDao,

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
	service.NewTalkListService,
	service.NewTalkMessageForwardService,
	service.NewEmoticonService,
	service.NewTalkRecordsService,
	service.NewContactService,
	service.NewContactsApplyService,

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
	v1.NewDownloadHandler,
	v1.NewEmoticonHandler,
	v1.NewTalkRecordsHandler,
	open.NewIndexHandler,
	ws.NewDefaultWebSocket,

	wire.Struct(new(handler.Handler), "*"),
	wire.Struct(new(provider.Services), "*"),

	// 持久化协程相关
	process.NewWsSubscribe,
	process.NewServerRun,
	process.NewProcessManage,
	process.NewImHeartbeat,
	handle.NewSubscribeConsume,
)

func Initialize(ctx context.Context) *provider.Services {
	panic(wire.Build(providerSet))
}
