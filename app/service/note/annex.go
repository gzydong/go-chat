package note

import (
	"context"
	"go-chat/app/dao/note"
	"go-chat/app/model"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service"
)

type ArticleAnnexService struct {
	*service.BaseService
	dao *note.ArticleAnnexDao
}

func NewArticleAnnexService(baseService *service.BaseService, dao *note.ArticleAnnexDao) *ArticleAnnexService {
	return &ArticleAnnexService{BaseService: baseService, dao: dao}
}

func (s *ArticleAnnexService) Dao() *note.ArticleAnnexDao {
	return s.dao
}

func (s *ArticleAnnexService) Create(ctx context.Context, data *model.ArticleAnnex) error {
	return s.Db().Create(data).Error
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
