package web

import (
	v1 "go-chat/internal/http/internal/handler/web/v1"
	"go-chat/internal/http/internal/handler/web/v1/article"
	"go-chat/internal/http/internal/handler/web/v1/contact"
	"go-chat/internal/http/internal/handler/web/v1/group"
	"go-chat/internal/http/internal/handler/web/v1/talk"
)

type V1 struct {
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
	GroupNotice   *group.Notice
	GroupApply    *group.Apply
	Contact       *contact.Contact
	ContactsApply *contact.Apply
	Article       *article.Article
	ArticleAnnex  *article.Annex
	ArticleClass  *article.Class
	ArticleTag    *article.Tag
	Test          *v1.Test
}

type Handler struct {
	V1 *V1
}
