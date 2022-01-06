package dao

import "go-chat/internal/model"

type SplitUploadDao struct {
	*BaseDao
}

func NewFileSplitUploadDao(baseDao *BaseDao) *SplitUploadDao {
	return &SplitUploadDao{BaseDao: baseDao}
}

func (dao *SplitUploadDao) GetSplitList(uploadId string) ([]*model.SplitUpload, error) {
	items := make([]*model.SplitUpload, 0)
	err := dao.Db().Model(&model.SplitUpload{}).Where("upload_id = ? and type = 2", uploadId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (dao *SplitUploadDao) GetFile(uid int, uploadId string) (*model.SplitUpload, error) {
	item := &model.SplitUpload{}

	err := dao.Db().First(item, "user_id = ? and upload_id = ? and type = 1", uid, uploadId).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}
