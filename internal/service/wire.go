package service

import (
	"github.com/google/wire"
	"go-chat/internal/service/note"
	"go-chat/internal/service/organize"
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

	note.NewArticleService,
	note.NewArticleTagService,
	note.NewArticleClassService,
	note.NewArticleAnnexService,
	organize.NewOrganizeDeptService,
	organize.NewOrganizeService,
	organize.NewPositionService,
	NewTemplateService,
	NewAuthService,
)
