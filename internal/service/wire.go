package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserService,

	wire.Struct(new(SmsService), "*"),
	wire.Bind(new(ISmsService), new(*SmsService)),

	NewTalkService,
	NewGroupService,
	NewGroupMemberService,
	NewGroupNoticeService,
	NewGroupApplyService,
	NewTalkSessionService,
	NewEmoticonService,
	NewTalkRecordsService,
	NewContactService,
	NewContactApplyService,
	NewContactGroupService,
	NewSplitUploadService,
	NewIpAddressService,

	wire.Struct(new(MessageService), "*"),
	wire.Bind(new(IMessageService), new(*MessageService)),

	wire.Struct(new(ArticleService), "*"),
	wire.Bind(new(IArticleService), new(*ArticleService)),

	NewArticleTagService,
	NewArticleClassService,
	NewArticleAnnexService,
	NewOrganizeDeptService,
	NewOrganizeService,
	NewPositionService,
	NewTemplateService,
	NewAuthService,
)
