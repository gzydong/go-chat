package repo

import (
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Emoticon struct {
	core.Repo[model.Emoticon]
}

func NewEmoticon(db *gorm.DB) *Emoticon {
	return &Emoticon{Repo: core.NewRepo[model.Emoticon](db)}
}

// GetUserInstallIds 获取用户激活的表情包
func (e *Emoticon) GetUserInstallIds(uid int) []int {
	var data model.UsersEmoticon
	if err := e.Repo.Db.First(&data, "user_id = ?", uid).Error; err != nil {
		return []int{}
	}

	return sliceutil.ParseIds(data.EmoticonIds)
}

// GetCustomizeList 获取自定义表情包
func (e *Emoticon) GetCustomizeList(uid int) ([]*model.EmoticonItem, error) {
	var items []*model.EmoticonItem
	if err := e.Repo.Db.Model(model.EmoticonItem{}).Where("emoticon_id = ? and user_id = ? order by id desc", 0, uid).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
