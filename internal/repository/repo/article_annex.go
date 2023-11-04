package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ArticleAnnex struct {
	ichat.Repo[model.ArticleAnnex]
}

func NewArticleAnnex(db *gorm.DB) *ArticleAnnex {
	return &ArticleAnnex{Repo: ichat.NewRepo[model.ArticleAnnex](db)}
}

func (a *ArticleAnnex) AnnexList(ctx context.Context, uid int, articleId int) ([]*model.ArticleAnnex, error) {
	return a.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("user_id = ? and article_id = ? and status = 1", uid, articleId)
	})
}

func (a *ArticleAnnex) RecoverList(ctx context.Context, uid int) ([]*model.RecoverAnnexItem, error) {

	fields := []string{
		"article_annex.id",
		"article_annex.article_id",
		"article.title",
		"article_annex.original_name",
		"article_annex.deleted_at",
	}

	query := a.Repo.Db.WithContext(ctx).Model(&model.ArticleAnnex{})
	query.Joins("left join article on article.id = article_annex.article_id")
	query.Where("article_annex.user_id = ? and article_annex.status = ?", uid, 2)

	items := make([]*model.RecoverAnnexItem, 0)
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
