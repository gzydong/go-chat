package repo

import (
	"go-chat/internal/repository/model"
)

type SplitUpload struct {
	*Base
}

func NewFileSplitUpload(base *Base) *SplitUpload {
	return &SplitUpload{Base: base}
}

func (s *SplitUpload) GetSplitList(uploadId string) ([]*model.SplitUpload, error) {
	items := make([]*model.SplitUpload, 0)
	err := s.Db.Model(&model.SplitUpload{}).Where("upload_id = ? and type = 2", uploadId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *SplitUpload) GetFile(uid int, uploadId string) (*model.SplitUpload, error) {
	item := &model.SplitUpload{}

	err := s.Db.First(item, "user_id = ? and upload_id = ? and type = 1", uid, uploadId).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}
