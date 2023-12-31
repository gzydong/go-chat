package service

import (
	"context"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type IArticleAnnexService interface {
	Create(ctx context.Context, data *model.ArticleAnnex) error
	UpdateStatus(ctx context.Context, uid int, id int, status int) error
	ForeverDelete(ctx context.Context, uid int, id int) error
}

type ArticleAnnexService struct {
	*repo.Source
	ArticleAnnex *repo.ArticleAnnex
	FileSystem   filesystem.IFilesystem
}

func (s *ArticleAnnexService) Create(ctx context.Context, data *model.ArticleAnnex) error {
	return s.ArticleAnnex.Create(ctx, data)
}

// UpdateStatus 更新附件状态
func (s *ArticleAnnexService) UpdateStatus(ctx context.Context, uid int, id int, status int) error {

	data := map[string]any{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	_, err := s.ArticleAnnex.UpdateWhere(ctx, data, "id = ? and user_id = ?", id, uid)
	return err
}

// ForeverDelete 永久删除笔记附件
func (s *ArticleAnnexService) ForeverDelete(ctx context.Context, uid int, id int) error {

	annex, err := s.ArticleAnnex.FindByWhere(ctx, "id = ? and user_id = ?", id, uid)
	if err != nil {
		return err
	}

	_ = s.FileSystem.Delete(s.FileSystem.BucketPrivateName(), annex.Path)

	return s.Source.Db().Delete(&model.ArticleAnnex{}, id).Error
}
