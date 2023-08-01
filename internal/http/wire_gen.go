// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/http/internal/handler/admin"
	v1_2 "go-chat/internal/http/internal/handler/admin/v1"
	"go-chat/internal/http/internal/handler/open"
	v1_3 "go-chat/internal/http/internal/handler/open/v1"
	"go-chat/internal/http/internal/handler/web"
	"go-chat/internal/http/internal/handler/web/v1"
	"go-chat/internal/http/internal/handler/web/v1/article"
	"go-chat/internal/http/internal/handler/web/v1/contact"
	"go-chat/internal/http/internal/handler/web/v1/group"
	"go-chat/internal/http/internal/handler/web/v1/talk"
	"go-chat/internal/http/internal/router"
	"go-chat/internal/logic"
	"go-chat/internal/provider"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/repo"
	"go-chat/internal/repository/repo/note"
	"go-chat/internal/repository/repo/organize"
	"go-chat/internal/service"
	note2 "go-chat/internal/service/note"
	organize2 "go-chat/internal/service/organize"
)

// Injectors from wire.go:

func Initialize(conf *config.Config) *AppProvider {
	client := provider.NewRedisClient(conf)
	smsStorage := cache.NewSmsStorage(client)
	smsService := service.NewSmsService(smsStorage)
	db := provider.NewMySQLClient(conf)
	users := repo.NewUsers(db)
	userService := service.NewUserService(users)
	common := v1.NewCommon(conf, smsService, userService)
	jwtTokenStorage := cache.NewTokenSessionStorage(client)
	redisLock := cache.NewRedisLock(client)
	source := repo.NewSource(db, client)
	httpClient := provider.NewHttpClient()
	requestClient := provider.NewRequestClient(httpClient)
	ipAddressService := service.NewIpAddressService(source, conf, requestClient)
	talkSession := repo.NewTalkSession(db)
	talkSessionService := service.NewTalkSessionService(source, talkSession)
	articleClass := note.NewArticleClass(db)
	articleClassService := note2.NewArticleClassService(source, articleClass)
	robot := repo.NewRobot(db)
	sequence := cache.NewSequence(client)
	repoSequence := repo.NewSequence(db, sequence)
	messageForwardLogic := logic.NewMessageForwardLogic(db, repoSequence)
	relation := cache.NewRelation(client)
	groupMember := repo.NewGroupMember(db, relation)
	splitUpload := repo.NewFileSplitUpload(db)
	vote := cache.NewVote(client)
	talkRecordsVote := repo.NewTalkRecordsVote(db, vote)
	filesystem := provider.NewFilesystem(conf)
	unreadStorage := cache.NewUnreadStorage(client)
	messageStorage := cache.NewMessageStorage(client)
	serverStorage := cache.NewSidStorage(client)
	clientStorage := cache.NewClientStorage(client, conf, serverStorage)
	messageService := service.NewMessageService(source, messageForwardLogic, groupMember, splitUpload, talkRecordsVote, filesystem, unreadStorage, messageStorage, serverStorage, clientStorage, repoSequence, robot)
	auth := v1.NewAuth(conf, userService, smsService, jwtTokenStorage, redisLock, ipAddressService, talkSessionService, articleClassService, robot, messageService)
	organizeOrganize := organize.NewOrganize(db)
	organizeService := organize2.NewOrganizeService(source, organizeOrganize)
	user := v1.NewUser(userService, smsService, organizeService)
	department := organize.NewDepartment(db)
	deptService := organize2.NewOrganizeDeptService(source, department)
	position := organize.NewPosition(db)
	positionService := organize2.NewPositionService(source, position)
	v1Organize := v1.NewOrganize(deptService, organizeService, positionService)
	talkService := service.NewTalkService(source, groupMember)
	contactRemark := cache.NewContactRemark(client)
	repoContact := repo.NewContact(db, contactRemark, relation)
	contactService := service.NewContactService(source, repoContact)
	repoGroup := repo.NewGroup(db)
	groupService := service.NewGroupService(source, repoGroup, groupMember, relation, repoSequence)
	authService := service.NewAuthService(organizeOrganize, repoContact, repoGroup, groupMember)
	session := talk.NewSession(talkService, talkSessionService, redisLock, userService, clientStorage, messageStorage, contactService, unreadStorage, contactRemark, groupService, authService)
	message := talk.NewMessage(talkService, authService, messageService, filesystem)
	talkRecords := repo.NewTalkRecords(db)
	talkRecordsService := service.NewTalkRecordsService(source, vote, talkRecordsVote, groupMember, talkRecords)
	groupMemberService := service.NewGroupMemberService(source, groupMember)
	records := talk.NewRecords(talkRecordsService, groupMemberService, filesystem, authService)
	emoticon := repo.NewEmoticon(db)
	emoticonService := service.NewEmoticonService(source, emoticon, filesystem)
	v1Emoticon := v1.NewEmoticon(filesystem, emoticonService, redisLock)
	splitUploadService := service.NewSplitUploadService(source, splitUpload, conf, filesystem)
	upload := v1.NewUpload(conf, filesystem, splitUploadService)
	groupNotice := repo.NewGroupNotice(db)
	groupNoticeService := service.NewGroupNoticeService(source, groupNotice)
	groupGroup := group.NewGroup(groupService, groupMemberService, talkSessionService, userService, redisLock, contactService, groupNoticeService, messageService)
	notice := group.NewNotice(groupNoticeService, groupMemberService, messageService)
	groupApply := repo.NewGroupApply(db)
	groupApplyService := service.NewGroupApplyService(source, groupApply)
	groupApplyStorage := cache.NewGroupApplyStorage(client)
	apply := group.NewApply(groupApplyService, groupMemberService, groupService, groupApplyStorage, client)
	contactContact := contact.NewContact(contactService, clientStorage, userService, talkSessionService, organizeService, messageService)
	contactApplyService := service.NewContactApplyService(source)
	contactApply := contact.NewApply(contactApplyService, userService, contactService, messageService)
	contactGroup := repo.NewContactGroup(db)
	contactGroupService := service.NewContactGroupService(source, contactGroup)
	group2 := contact.NewGroup(contactGroupService, contactService)
	articleService := note2.NewArticleService(source)
	articleAnnex := note.NewArticleAnnex(db)
	articleAnnexService := note2.NewArticleAnnexService(source, articleAnnex, filesystem)
	articleArticle := article.NewArticle(articleService, filesystem, articleAnnexService)
	annex := article.NewAnnex(articleAnnexService, filesystem)
	class := article.NewClass(articleClassService)
	articleTagService := note2.NewArticleTagService(source)
	tag := article.NewTag(articleTagService)
	publish := talk.NewPublish(authService, messageService)
	webV1 := &web.V1{
		Common:       common,
		Auth:         auth,
		User:         user,
		Organize:     v1Organize,
		Talk:         session,
		TalkMessage:  message,
		TalkRecords:  records,
		Emoticon:     v1Emoticon,
		Upload:       upload,
		Group:        groupGroup,
		GroupNotice:  notice,
		GroupApply:   apply,
		Contact:      contactContact,
		ContactApply: contactApply,
		ContactGroup: group2,
		Article:      articleArticle,
		ArticleAnnex: annex,
		ArticleClass: class,
		ArticleTag:   tag,
		Message:      publish,
	}
	webHandler := &web.Handler{
		V1: webV1,
	}
	index := v1_2.NewIndex()
	captchaStorage := cache.NewCaptchaStorage(client)
	repoAdmin := repo.NewAdmin(db)
	v1Auth := v1_2.NewAuth(conf, captchaStorage, repoAdmin, jwtTokenStorage)
	adminV1 := &admin.V1{
		Index: index,
		Auth:  v1Auth,
	}
	v2 := &admin.V2{}
	adminHandler := &admin.Handler{
		V1: adminV1,
		V2: v2,
	}
	v1Index := v1_3.NewIndex()
	openV1 := &open.V1{
		Index: v1Index,
	}
	openHandler := &open.Handler{
		V1: openV1,
	}
	handlerHandler := &handler.Handler{
		Api:   webHandler,
		Admin: adminHandler,
		Open:  openHandler,
	}
	engine := router.NewRouter(conf, handlerHandler, jwtTokenStorage)
	appProvider := &AppProvider{
		Config: conf,
		Engine: engine,
	}
	return appProvider
}

// wire.go:

var providerSet = wire.NewSet(provider.NewMySQLClient, provider.NewRedisClient, provider.NewHttpClient, provider.NewEmailClient, provider.NewFilesystem, provider.NewRequestClient, router.NewRouter, wire.Struct(new(web.Handler), "*"), wire.Struct(new(admin.Handler), "*"), wire.Struct(new(open.Handler), "*"), wire.Struct(new(handler.Handler), "*"), wire.Struct(new(AppProvider), "*"))

var daoProviderSet = wire.NewSet(repo.NewSource, repo.NewContact, repo.NewContactGroup, repo.NewGroupMember, repo.NewUsers, repo.NewGroup, repo.NewGroupApply, repo.NewTalkRecords, repo.NewGroupNotice, repo.NewTalkSession, repo.NewEmoticon, repo.NewTalkRecordsVote, repo.NewFileSplitUpload, note.NewArticleClass, note.NewArticleAnnex, organize.NewDepartment, organize.NewOrganize, organize.NewPosition, repo.NewRobot, repo.NewSequence, repo.NewAdmin)

var serviceProviderSet = wire.NewSet(service.NewUserService, service.NewSmsService, service.NewTalkService, service.NewGroupService, service.NewGroupMemberService, service.NewGroupNoticeService, service.NewGroupApplyService, service.NewTalkSessionService, service.NewEmoticonService, service.NewTalkRecordsService, service.NewContactService, service.NewContactApplyService, service.NewContactGroupService, service.NewSplitUploadService, service.NewIpAddressService, service.NewMessageService, note2.NewArticleService, note2.NewArticleTagService, note2.NewArticleClassService, note2.NewArticleAnnexService, organize2.NewOrganizeDeptService, organize2.NewOrganizeService, organize2.NewPositionService, service.NewTemplateService, service.NewAuthService, logic.NewMessageForwardLogic)
