package service

import (
	"context"
	"errors"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ IArticleClassService = (*ArticleClassService)(nil)

type IArticleClassService interface {
	List(ctx context.Context, uid int) ([]*model.ArticleClassItem, error)
	Create(ctx context.Context, uid int, name string) (int, error)
	Update(ctx context.Context, uid, cid int, name string) error
	Delete(ctx context.Context, uid, cid int) error
	Sort(ctx context.Context, uid, cid, mode int) error
}

type ArticleClassService struct {
	*repo.Source
	ArticleClass *repo.ArticleClass
}

// List 分类列表
func (s *ArticleClassService) List(ctx context.Context, uid int) ([]*model.ArticleClassItem, error) {
	items := make([]*model.ArticleClassItem, 0)

	err := s.ArticleClass.Model(ctx).Select("id", "class_name", "is_default").Where("user_id = ?", uid).Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	data, err := s.ArticleClass.GroupCount(uid)
	if err != nil {
		return nil, err
	}

	items = append(items, &model.ArticleClassItem{
		ClassName: "默认分类",
	})

	for i := range items {
		if num, ok := data[items[i].Id]; ok {
			items[i].Count = num
		}
	}

	return items, nil
}

// Create 创建分类
func (s *ArticleClassService) Create(ctx context.Context, uid int, name string) (int, error) {
	data := &model.ArticleClass{
		UserId:    uid,
		ClassName: name,
		Sort:      1,
	}

	err := s.Source.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("article_class").Where("user_id = ?", uid).Updates(map[string]any{
			"sort": gorm.Expr("sort + 1"),
		}).Error; err != nil {
			return err
		}

		return tx.Create(data).Error
	})
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (s *ArticleClassService) Update(ctx context.Context, uid, cid int, name string) error {
	_, err := s.ArticleClass.UpdateWhere(ctx, map[string]any{"class_name": name}, "id = ? and user_id = ?", cid, uid)
	return err
}

func (s *ArticleClassService) Delete(ctx context.Context, uid, cid int) error {

	var num int64
	if err := s.Source.Db().WithContext(ctx).Table("article").Where("user_id = ? and class_id = ?", uid, cid).Count(&num).Error; err != nil {
		return err
	}

	if num > 0 {
		return errors.New("分类已被使用不能删除")
	}

	return s.Source.Db().Delete(&model.ArticleClass{}, "id = ? and user_id = ? and is_default = 0", cid, uid).Error
}

func (s *ArticleClassService) Sort(ctx context.Context, uid, cid, mode int) error {

	item, err := s.ArticleClass.FindByWhere(ctx, "id = ? and user_id = ?", cid, uid)
	if err != nil {
		return err
	}

	if mode == 1 {
		maxSort, err := s.ArticleClass.MaxSort(ctx, uid)
		if err != nil {
			return err
		}

		if maxSort == item.Sort {
			return nil
		}

		return s.Source.Db().Transaction(func(tx *gorm.DB) error {
			if err := tx.Table("article_class").Where("user_id = ? and sort = ?", uid, item.Sort+1).Updates(map[string]any{
				"sort": gorm.Expr("sort - 1"),
			}).Error; err != nil {
				return err
			}

			if err := tx.Table("article_class").Where("id = ? and user_id = ?", cid, uid).Updates(map[string]any{
				"sort": gorm.Expr("sort + 1"),
			}).Error; err != nil {
				return err
			}

			return nil
		})
	} else {
		minSort, err := s.ArticleClass.MinSort(ctx, uid)
		if err != nil {
			return err
		}

		if minSort == item.Sort {
			return nil
		}

		return s.Source.Db().Transaction(func(tx *gorm.DB) error {
			if err := tx.Table("article_class").Where("user_id = ? and sort = ?", uid, item.Sort-1).Updates(map[string]any{
				"sort": gorm.Expr("sort + 1"),
			}).Error; err != nil {
				return err
			}

			if err := tx.Table("article_class").Where("id = ? and user_id = ?", cid, uid).Updates(map[string]any{
				"sort": gorm.Expr("sort - 1"),
			}).Error; err != nil {
				return err
			}

			return nil
		})
	}
}

// SetDefaultClass 设置默认分类
func (s *ArticleClassService) SetDefaultClass(ctx context.Context, uid int) {

	_, err := s.ArticleClass.QueryExist(ctx, "id = ? and is_default = ?", uid, 1)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	_ = s.ArticleClass.Create(ctx, &model.ArticleClass{
		UserId:    uid,
		ClassName: "默认分类",
		Sort:      1,
		IsDefault: 1,
	})
}
