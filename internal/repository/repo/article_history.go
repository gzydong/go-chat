package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ArticleHistory struct {
	core.Repo[model.ArticleHistory]
}

func NewArticleHistory(db *gorm.DB) *ArticleHistory {
	return &ArticleHistory{Repo: core.NewRepo[model.ArticleHistory](db)}
}
