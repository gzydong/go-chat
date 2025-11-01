package repo

import (
	"github.com/gzydong/go-chat/internal/pkg/core"
	"github.com/gzydong/go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Article struct {
	core.Repo[model.Article]
}

func NewArticle(db *gorm.DB) *Article {
	return &Article{Repo: core.NewRepo[model.Article](db)}
}
