package web

import (
	v12 "go-chat/internal/api/handler/web/v1"
	article2 "go-chat/internal/api/handler/web/v1/article"
	contact2 "go-chat/internal/api/handler/web/v1/contact"
	group2 "go-chat/internal/api/handler/web/v1/group"
	talk2 "go-chat/internal/api/handler/web/v1/talk"
)

type V1 struct {
	Common       *v12.Common
	Auth         *v12.Auth
	User         *v12.User
	Organize     *v12.Organize
	Talk         *talk2.Session
	TalkMessage  *talk2.Message
	TalkRecords  *talk2.Records
	Emoticon     *v12.Emoticon
	Upload       *v12.Upload
	Group        *group2.Group
	GroupNotice  *group2.Notice
	GroupApply   *group2.Apply
	Contact      *contact2.Contact
	ContactApply *contact2.Apply
	ContactGroup *contact2.Group
	Article      *article2.Article
	ArticleAnnex *article2.Annex
	ArticleClass *article2.Class
	ArticleTag   *article2.Tag
	Message      *talk2.Publish
}

type Handler struct {
	V1 *V1
}
