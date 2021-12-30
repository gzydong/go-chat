package handler

import (
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/api/v1/article"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
)

type Handler struct {
	Common           *v1.Common
	Auth             *v1.Auth
	User             *v1.User
	TalkMessage      *v1.TalkMessage
	Talk             *v1.Talk
	TalkRecords      *v1.TalkRecords
	Emoticon         *v1.Emoticon
	Upload           *v1.Upload
	Index            *open.Index
	DefaultWebSocket *ws.DefaultWebSocket
	Group            *v1.Group
	GroupNotice      *v1.GroupNotice
	Contact          *v1.Contact
	ContactsApply    *v1.ContactApply
	Article          *article.Article
	ArticleAnnex     *article.Annex
	ArticleClass     *article.Class
	ArticleTag       *article.Tag
}
