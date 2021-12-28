package note

import (
	"context"
	"go-chat/app/model"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service"
)

type ArticleAnnexService struct {
	*service.BaseService
}

func NewArticleAnnexService(baseService *service.BaseService) *ArticleAnnexService {
	return &ArticleAnnexService{BaseService: baseService}
}

func (s *ArticleAnnexService) Create(ctx context.Context, data *model.ArticleAnnex) error {
	return s.Db().Create(data).Error
}

func (s *ArticleAnnexService) AnnexList(ctx context.Context, uid int, articleId int) ([]*model.ArticleAnnex, error) {
	items := make([]*model.ArticleAnnex, 0)

	err := s.Db().Model(&model.ArticleAnnex{}).Where("user_id = ? and article_id = ? and status = 1", uid, articleId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ArticleAnnexService) FindById(ctx context.Context, id int) (*model.ArticleAnnex, error) {
	item := &model.ArticleAnnex{}

	if err := s.Db().First(item, id).Error; err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ArticleAnnexService) UpdateStatus(ctx context.Context, uid int, id int, status int) error {

	data := map[string]interface{}{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	return s.Db().Model(&model.ArticleAnnex{}).Where("id = ? and user_id = ?", id, uid).Updates(data).Error
}

func (s *ArticleAnnexService) RecoverList(ctx context.Context, uid int) ([]*model.RecoverAnnexItem, error) {

	fields := []string{
		"article_annex.id",
		"article_annex.article_id",
		"article.title",
		"article_annex.original_name",
		"article_annex.deleted_at",
	}

	query := s.Db().Model(&model.ArticleAnnex{})
	query.Joins("left join article on article.id = article_annex.article_id")
	query.Where("article_annex.user_id = ? and article_annex.status = ?", uid, 2)

	items := make([]*model.RecoverAnnexItem, 0)
	if err := query.Select(fields).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
