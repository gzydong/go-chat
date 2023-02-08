package web

import (
	"github.com/google/wire"
	"go-chat/internal/http/internal/handler/web/v1"
	"go-chat/internal/http/internal/handler/web/v1/article"
	"go-chat/internal/http/internal/handler/web/v1/contact"
	"go-chat/internal/http/internal/handler/web/v1/group"
	"go-chat/internal/http/internal/handler/web/v1/talk"
)

var ProviderSet = wire.NewSet(
	v1.NewAuth,
	v1.NewCommon,
	v1.NewUser,
	v1.NewOrganize,
	contact.NewContact,
	contact.NewApply,
	contact.NewGroup,
	group.NewGroup,
	group.NewApply,
	group.NewNotice,
	talk.NewSession,
	talk.NewMessage,
	v1.NewUpload,
	v1.NewEmoticon,
	talk.NewRecords,
	article.NewAnnex,
	article.NewArticle,
	article.NewClass,
	article.NewTag,
	talk.NewSendMessage,

	wire.Struct(new(V1), "*"),
)
