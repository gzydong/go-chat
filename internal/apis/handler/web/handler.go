package web

import (
	v1 "github.com/gzydong/go-chat/internal/apis/handler/web/v1"
	"github.com/gzydong/go-chat/internal/apis/handler/web/v1/article"
	"github.com/gzydong/go-chat/internal/apis/handler/web/v1/contact"
	"github.com/gzydong/go-chat/internal/apis/handler/web/v1/group"
	"github.com/gzydong/go-chat/internal/apis/handler/web/v1/talk"
	"github.com/gzydong/go-chat/internal/repository/repo"
)

type V1 struct {
	Common       *v1.Common
	Auth         *v1.Auth
	User         *v1.User
	Organize     *v1.Organize
	Talk         *talk.Session
	TalkMessage  *talk.Message
	Emoticon     *v1.Emoticon
	Upload       *v1.Upload
	Group        *group.Group
	GroupNotice  *group.Notice
	GroupApply   *group.Apply
	GroupVote    *group.Vote
	Contact      *contact.Contact
	ContactApply *contact.Apply
	ContactGroup *contact.Group
	Article      *article.Article
	ArticleAnnex *article.Annex
	ArticleClass *article.Class
	ArticleTag   *article.Tag
	Message      *talk.Publish
}

type Handler struct {
	V1       *V1
	UserRepo *repo.Users
}
