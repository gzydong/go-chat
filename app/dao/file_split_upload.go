package dao

import "go-chat/app/model"

type FileSplitUploadDao struct {
	*BaseDao
}

func NewFileSplitUploadDao(baseDao *BaseDao) *FileSplitUploadDao {
	return &FileSplitUploadDao{BaseDao: baseDao}
}

func (dao *FileSplitUploadDao) GetSplitList(uploadId string) ([]*model.FileSplitUpload, error) {
	items := make([]*model.FileSplitUpload, 0)
	err := dao.Db().Model(model.FileSplitUpload{}).Where("upload_id = ? and type = 2", uploadId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
