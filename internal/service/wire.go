package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserService,
	NewSmsService,
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

	NewArticleService,
	NewArticleTagService,
	NewArticleClassService,
	NewArticleAnnexService,
	NewOrganizeDeptService,
	NewOrganizeService,
	NewPositionService,
	NewTemplateService,
	NewAuthService,
)
