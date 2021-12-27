package note

import (
	"context"
	"go-chat/app/model"
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
