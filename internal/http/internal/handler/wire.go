package handler

import (
	"github.com/google/wire"
	v1 "go-chat/internal/http/internal/handler/api/v1"
	"go-chat/internal/http/internal/handler/api/v1/article"
	"go-chat/internal/http/internal/handler/api/v1/contact"
	"go-chat/internal/http/internal/handler/api/v1/group"
	"go-chat/internal/http/internal/handler/api/v1/talk"
)

var ProviderSet = wire.NewSet(
	v1.NewAuthHandler,
	v1.NewCommonHandler,
	v1.NewUserHandler,
	v1.NewOrganizeHandler,
	contact.NewContactHandler,
	contact.NewContactsApplyHandler,
	group.NewGroupHandler,
	group.NewApplyHandler,
	group.NewGroupNoticeHandler,
	talk.NewTalkHandler,
	talk.NewTalkMessageHandler,
	v1.NewUploadHandler,
	v1.NewEmoticonHandler,
	talk.NewTalkRecordsHandler,
	article.NewAnnexHandler,
	article.NewArticleHandler,
	article.NewClassHandler,
	article.NewTagHandler,
)
