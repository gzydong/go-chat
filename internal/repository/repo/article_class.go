package repo

import (
	"context"

	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ArticleClass struct {
	ichat.Repo[model.ArticleClass]
}

func NewArticleClass(db *gorm.DB) *ArticleClass {
	return &ArticleClass{Repo: ichat.NewRepo[model.ArticleClass](db)}
}

func (a *ArticleClass) MaxSort(ctx context.Context, uid int) (int, error) {
	var sort int
	err := a.Repo.Model(ctx).Select("max(sort)").Where("user_id = ?", uid).Scan(&sort).Error
	if err != nil {
		return 0, err
	}

	return sort, nil
}

func (a *ArticleClass) MinSort(ctx context.Context, uid int) (int, error) {
	var sort int
	err := a.Repo.Model(ctx).Select("min(sort)").Where("user_id = ?", uid).Scan(&sort).Error
	if err != nil {
		return 0, err
	}

	return sort, nil
}

type ClassGroupCount struct {
	ClassId int `json:"class_id"`
	Count   int `json:"count"`
}

func (a *ArticleClass) GroupCount(uid int) (map[int]int, error) {
	var items []ClassGroupCount
	if err := a.Repo.Db.Table("article").
		Select("class_id", "ifnull(count(*),0) as count").
		Where("user_id = ? and status = 1", uid).
		Group("class_id").
		Scan(&items).Error; err != nil {
		return nil, err
	}

	maps := make(map[int]int)
	for _, item := range items {
		maps[item.ClassId] = item.Count
	}

	return maps, nil
}
