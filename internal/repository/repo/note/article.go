package note

import (
	"go-chat/internal/repository/repo"
)

type Article struct {
	*repo.Base
}

func NewArticle(base *repo.Base) *Article {
	return &Article{Base: base}
}
