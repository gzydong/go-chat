package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type SplitUpload struct {
	ichat.Repo[model.SplitUpload]
}

func NewFileSplitUpload(db *gorm.DB) *SplitUpload {
	return &SplitUpload{Repo: ichat.NewRepo[model.SplitUpload](db)}
}

func (s *SplitUpload) GetSplitList(ctx context.Context, uploadId string) ([]*model.SplitUpload, error) {
	return s.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("upload_id = ? and type = 2", uploadId)
	})
}

func (s *SplitUpload) GetFile(ctx context.Context, uid int, uploadId string) (*model.SplitUpload, error) {
	return s.Repo.FindByWhere(ctx, "user_id = ? and upload_id = ? and type = 1", uid, uploadId)
}
