package repo

import (
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
)

type IEmoticon interface {
	IBase
	FindById(emoticonId int) (*model.Emoticon, error)
	GetUserInstallIds(uid int) []int
	GetSystemEmoticonList() ([]*model.Emoticon, error)
	GetDetailsAll(emoticonId int, uid int) ([]*model.EmoticonItem, error)
}

type Emoticon struct {
	*Base
}

func NewEmoticon(base *Base) *Emoticon {
	return &Emoticon{Base: base}
}

func (e *Emoticon) FindById(id int) (*model.Emoticon, error) {
	var data *model.Emoticon

	if err := e.Db.First(&data, id).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (e *Emoticon) FindByIds(ids []int) ([]*model.Emoticon, error) {

	items := make([]*model.Emoticon, 0)
	if err := e.Db.Find(&items, ids).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// GetUserInstallIds 获取用户激活的表情包
func (e *Emoticon) GetUserInstallIds(uid int) []int {
	data := &model.UsersEmoticon{}

	if err := e.Db.First(data, "user_id = ?", uid).Error; err != nil {
		return []int{}
	}

	return sliceutil.ParseIds(data.EmoticonIds)
}

// GetSystemEmoticonList 获取系统表情包分组列表
func (e *Emoticon) GetSystemEmoticonList() ([]*model.Emoticon, error) {
	var (
		err   error
		items = make([]*model.Emoticon, 0)
	)

	err = e.Db.Model(&model.Emoticon{}).Where("status = ?", 0).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetDetailsAll 获取系统表情包分组详情列表
func (e *Emoticon) GetDetailsAll(emoticonId, uid int) ([]*model.EmoticonItem, error) {
	var (
		err   error
		items = make([]*model.EmoticonItem, 0)
	)

	if err = e.Db.Model(&model.EmoticonItem{}).Where("emoticon_id = ? and user_id = ? order by id desc", emoticonId, uid).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
