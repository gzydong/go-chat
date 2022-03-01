package note

import (
	"context"

	"go-chat/internal/dao/note"
	"go-chat/internal/model"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/service"
)

type ArticleAnnexService struct {
	*service.BaseService
	dao        *note.ArticleAnnexDao
	fileSystem *filesystem.Filesystem
}

func NewArticleAnnexService(baseService *service.BaseService, dao *note.ArticleAnnexDao, fileSystem *filesystem.Filesystem) *ArticleAnnexService {
	return &ArticleAnnexService{BaseService: baseService, dao: dao, fileSystem: fileSystem}
}

func (s *ArticleAnnexService) Dao() *note.ArticleAnnexDao {
	return s.dao
}

func (s *ArticleAnnexService) Create(ctx context.Context, data *model.ArticleAnnex) error {
	return s.Db().Create(data).Error
}

// UpdateStatus 更新附件状态
func (s *ArticleAnnexService) UpdateStatus(ctx context.Context, uid int, id int, status int) error {

	data := map[string]interface{}{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	return s.Db().Model(&model.ArticleAnnex{}).Where("id = ? and user_id = ?", id, uid).Updates(data).Error
}

// ForeverDelete 永久删除笔记附件
func (s *ArticleAnnexService) ForeverDelete(ctx context.Context, uid int, id int) error {
	var annex *model.ArticleAnnex

	if err := s.Db().First(&annex, "id = ? and user_id = ?", id, uid).Error; err != nil {
		return err
	}

	switch annex.Drive {
	case 1:
		_ = s.fileSystem.Local.Delete(annex.Path)
	case 2:
		_ = s.fileSystem.Cos.Delete(annex.Path)
	}

	return s.Db().Delete(&model.ArticleAnnex{}, id).Error
}
