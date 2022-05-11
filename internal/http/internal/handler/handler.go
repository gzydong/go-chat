package handler

import (
	"go-chat/internal/http/internal/handler/api/v1"
	"go-chat/internal/http/internal/handler/api/v1/article"
	"go-chat/internal/http/internal/handler/api/v1/contact"
	"go-chat/internal/http/internal/handler/api/v1/group"
	"go-chat/internal/http/internal/handler/api/v1/talk"
)

type Handler struct {
	Common        *v1.Common
	Auth          *v1.Auth
	User          *v1.User
	Organize      *v1.Organize
	TalkMessage   *talk.Message
	Talk          *talk.Talk
	TalkRecords   *talk.Records
	Emoticon      *v1.Emoticon
	Upload        *v1.Upload
	Group         *group.Group
	GroupNotice   *group.GroupNotice
	Contact       *contact.Contact
	ContactsApply *contact.ContactApply
	Article       *article.Article
	ArticleAnnex  *article.Annex
	ArticleClass  *article.Class
	ArticleTag    *article.Tag
}
