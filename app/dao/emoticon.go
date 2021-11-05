package dao

import (
	"go-chat/app/model"
)

type EmoticonDao struct {
	*Base
}

func NewEmoticonDao(base *Base) *EmoticonDao {
	return &EmoticonDao{Base: base}
}

func (dao *EmoticonDao) GetSystemEmoticonList() ([]*model.Emoticon, error) {
	var (
		err   error
		items []*model.Emoticon
	)

	err = dao.Db.Model(model.Emoticon{}).Where("status = ?", 0).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
