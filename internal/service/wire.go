package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	wire.Struct(new(UserService), "*"),
	wire.Bind(new(IUserService), new(*UserService)),

	wire.Struct(new(SmsService), "*"),
	wire.Bind(new(ISmsService), new(*SmsService)),

	wire.Struct(new(TalkService), "*"),
	wire.Bind(new(ITalkService), new(*TalkService)),

	wire.Struct(new(GroupService), "*"),
	wire.Bind(new(IGroupService), new(*GroupService)),

	wire.Struct(new(GroupMemberService), "*"),
	wire.Bind(new(IGroupMemberService), new(*GroupMemberService)),

	wire.Struct(new(GroupNoticeService), "*"),
	wire.Bind(new(IGroupNoticeService), new(*GroupNoticeService)),

	wire.Struct(new(GroupApplyService), "*"),
	wire.Bind(new(IGroupApplyService), new(*GroupApplyService)),

	wire.Struct(new(TalkSessionService), "*"),
	wire.Bind(new(ITalkSessionService), new(*TalkSessionService)),

	wire.Struct(new(EmoticonService), "*"),
	wire.Bind(new(IEmoticonService), new(*EmoticonService)),

	wire.Struct(new(TalkRecordsService), "*"),
	wire.Bind(new(ITalkRecordsService), new(*TalkRecordsService)),

	wire.Struct(new(ContactService), "*"),
	wire.Bind(new(IContactService), new(*ContactService)),

	wire.Struct(new(ContactApplyService), "*"),
	wire.Bind(new(IContactApplyService), new(*ContactApplyService)),

	wire.Struct(new(ContactGroupService), "*"),
	wire.Bind(new(IContactGroupService), new(*ContactGroupService)),

	wire.Struct(new(SplitUploadService), "*"),
	wire.Bind(new(ISplitUploadService), new(*SplitUploadService)),

	wire.Struct(new(IpAddressService), "*"),
	wire.Bind(new(IIpAddressService), new(*IpAddressService)),

	wire.Struct(new(MessageService), "*"),
	wire.Bind(new(IMessageService), new(*MessageService)),

	wire.Struct(new(ArticleService), "*"),
	wire.Bind(new(IArticleService), new(*ArticleService)),

	wire.Struct(new(ArticleTagService), "*"),
	wire.Bind(new(IArticleTagService), new(*ArticleTagService)),

	wire.Struct(new(ArticleClassService), "*"),
	wire.Bind(new(IArticleClassService), new(*ArticleClassService)),

	wire.Struct(new(ArticleAnnexService), "*"),
	wire.Bind(new(IArticleAnnexService), new(*ArticleAnnexService)),

	wire.Struct(new(TemplateService), "*"),
	wire.Bind(new(ITemplateService), new(*TemplateService)),

	wire.Struct(new(AuthService), "*"),
	wire.Bind(new(IAuthService), new(*AuthService)),
)
