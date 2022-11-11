package note

import (
	"context"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type ArticleAnnex struct {
	*repo.Base
}

func NewArticleAnnex(baseDao *repo.Base) *ArticleAnnex {
	return &ArticleAnnex{Base: baseDao}
}

func (repo *ArticleAnnex) FindById(ctx context.Context, id int) (*model.ArticleAnnex, error) {
	item := &model.ArticleAnnex{}

	if err := repo.Db().First(item, id).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (repo *ArticleAnnex) AnnexList(ctx context.Context, uid int, articleId int) ([]*model.ArticleAnnex, error) {
	items := make([]*model.ArticleAnnex, 0)

	err := repo.Db().Model(&model.ArticleAnnex{}).Where("user_id = ? and article_id = ? and status = 1", uid, articleId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *ArticleAnnex) RecoverList(ctx context.Context, uid int) ([]*model.RecoverAnnexItem, error) {

	fields := []string{
		"article_annex.id",
		"article_annex.article_id",
		"article.title",
		"article_annex.original_name",
		"article_annex.deleted_at",
	}

	query := repo.Db().Model(&model.ArticleAnnex{})
	query.Joins("left join article on article.id = article_annex.article_id")
	query.Where("article_annex.user_id = ? and article_annex.status = ?", uid, 2)

	items := make([]*model.RecoverAnnexItem, 0)
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
