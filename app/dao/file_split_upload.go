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
	err := dao.Db().Model(&model.FileSplitUpload{}).Where("upload_id = ? and type = 2", uploadId).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (dao *FileSplitUploadDao) GetFile(uid int, uploadId string) (*model.FileSplitUpload, error) {
	item := &model.FileSplitUpload{}

	err := dao.Db().First(item, "user_id = ? and upload_id = ? and type = 1", uid, uploadId).Error
	if err != nil {
		return nil, err
	}

	return item, nil
}
