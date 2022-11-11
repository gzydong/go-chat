package repo

import (
	"go-chat/internal/repository/model"
)

type SplitUpload struct {
	*Base
}

func NewFileSplitUpload(baseDao *Base) *SplitUpload {
	return &SplitUpload{Base: baseDao}
}

func (repo *SplitUpload) GetSplitList(uploadId string) ([]*model.SplitUpload, error) {
	items := make([]*model.SplitUpload, 0)
	err := repo.db.Model(&model.SplitUpload{}).Where("upload_id = ? and type = 2", uploadId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SplitUpload) GetFile(uid int, uploadId string) (*model.SplitUpload, error) {
	item := &model.SplitUpload{}

	err := repo.db.First(item, "user_id = ? and upload_id = ? and type = 1", uid, uploadId).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}
