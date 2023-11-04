package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type Emoticon struct {
	ichat.Repo[model.Emoticon]
}

func NewEmoticon(db *gorm.DB) *Emoticon {
	return &Emoticon{Repo: ichat.NewRepo[model.Emoticon](db)}
}

// GetUserInstallIds 获取用户激活的表情包
func (e *Emoticon) GetUserInstallIds(uid int) []int {
	var data model.UsersEmoticon
	if err := e.Repo.Db.First(&data, "user_id = ?", uid).Error; err != nil {
		return []int{}
	}

	return sliceutil.ParseIds(data.EmoticonIds)
}

// GetSystemEmoticonList 获取系统表情包分组列表
func (e *Emoticon) GetSystemEmoticonList(ctx context.Context) ([]*model.Emoticon, error) {
	return e.Repo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("status = ?", 0)
	})
}

// GetDetailsAll 获取系统表情包分组详情列表
func (e *Emoticon) GetDetailsAll(emoticonId, uid int) ([]*model.EmoticonItem, error) {
	var items []*model.EmoticonItem
	if err := e.Repo.Db.Model(model.EmoticonItem{}).Where("emoticon_id = ? and user_id = ? order by id desc", emoticonId, uid).Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
