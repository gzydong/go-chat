package dao

import (
	"go-chat/app/model"
	"go-chat/app/pkg/slice"
)

type EmoticonDao struct {
	*BaseDao
}

func NewEmoticonDao(base *BaseDao) *EmoticonDao {
	return &EmoticonDao{BaseDao: base}
}

func (dao *EmoticonDao) FindById(emoticonId int) (*model.Emoticon, error) {
	var data *model.Emoticon

	if err := dao.Db.First(&data, emoticonId).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// GetUserInstallIds 获取用户激活的表情包
func (dao *EmoticonDao) GetUserInstallIds(uid int) []int {
	data := &model.UsersEmoticon{}

	if err := dao.Db.First(data, "user_id = ?", uid).Error; err != nil {
		return []int{}
	}

	return slice.ParseIds(data.EmoticonIds)
}

// GetSystemEmoticonList 获取系统表情包分组列表
func (dao *EmoticonDao) GetSystemEmoticonList() ([]*model.Emoticon, error) {
	var (
		err   error
		items = make([]*model.Emoticon, 0)
	)

	err = dao.Db.Model(&model.Emoticon{}).Where("status = ?", 0).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetDetailsAll 获取系统表情包分组详情列表
func (dao *EmoticonDao) GetDetailsAll(emoticonId, uid int) ([]*model.EmoticonItem, error) {
	var (
		err   error
		items = make([]*model.EmoticonItem, 0)
	)

	if err = dao.Db.Model(&model.EmoticonItem{}).Where("emoticon_id = ? and user_id = ?", emoticonId, uid).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
