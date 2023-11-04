package repo

import (
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Article struct {
	ichat.Repo[model.Article]
}

func NewArticle(db *gorm.DB) *Article {
	return &Article{Repo: ichat.NewRepo[model.Article](db)}
}
