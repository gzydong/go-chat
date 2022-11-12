package note

import (
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
)

type ArticleClass struct {
	*repo.Base
}

func NewArticleClass(base *repo.Base) *ArticleClass {
	return &ArticleClass{Base: base}
}

func (repo *ArticleClass) MaxSort(uid int) (int, error) {
	var sort int

	err := repo.Db.Model(&model.ArticleClass{}).Select("max(sort)").Where("user_id = ?", uid).Scan(&sort).Error
	if err != nil {
		return 0, err
	}

	return sort, nil
}

func (repo *ArticleClass) MinSort(uid int) (int, error) {
	var sort int

	err := repo.Db.Model(&model.ArticleClass{}).Select("min(sort)").Where("user_id = ?", uid).Scan(&sort).Error
	if err != nil {
		return 0, err
	}

	return sort, nil
}

type ClassCount struct {
	ClassId int `json:"class_id"`
	Count   int `json:"count"`
}

func (repo *ArticleClass) GroupCount(uid int) (map[int]int, error) {
	items := make([]*ClassCount, 0)
	if err := repo.Db.Model(&model.Article{}).Select("class_id", "count(*) as count").Where("user_id = ? and status = 1", uid).Group("class_id").Scan(&items).Error; err != nil {
		return nil, err
	}

	maps := make(map[int]int)

	for _, item := range items {
		maps[item.ClassId] = item.Count
	}

	return maps, nil
}
