package repo

import (
	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ArticleHistory struct {
	core.Repo[model.ArticleHistory]
}

func NewArticleHistory(db *gorm.DB) *ArticleHistory {
	return &ArticleHistory{Repo: core.NewRepo[model.ArticleHistory](db)}
}
