package note

import (
	"context"
	"errors"

	"go-chat/internal/repository/dao/note"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"

	"go-chat/internal/service"
)

type ArticleClassService struct {
	*service.BaseService
	dao *note.ArticleClassDao
}

func NewArticleClassService(baseService *service.BaseService, dao *note.ArticleClassDao) *ArticleClassService {
	return &ArticleClassService{BaseService: baseService, dao: dao}
}

// List 分类列表
func (s *ArticleClassService) List(ctx context.Context, uid int) ([]*model.ArticleClassItem, error) {
	items := make([]*model.ArticleClassItem, 0)

	err := s.Db().Model(&model.ArticleClass{}).Select("id", "class_name", "is_default").Where("user_id = ?", uid).Order("sort asc").Scan(&items).Error
	if err != nil {
		return nil, err
	}

	data, err := s.dao.GroupCount(uid)
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

	err := s.Db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.ArticleClass{}).Where("user_id = ?", uid).Updates(map[string]interface{}{
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
	return s.Db().Model(&model.ArticleClass{}).Where("id = ? and user_id = ?", cid, uid).UpdateColumn("class_name", name).Error
}

func (s *ArticleClassService) Delete(ctx context.Context, uid, cid int) error {
	var num int64
	if err := s.Db().Model(&model.Article{}).Where("user_id = ? and class_id = ?", uid, cid).Count(&num).Error; err != nil {
		return err
	}

	if num > 0 {
		return errors.New("分类已被使用不能删除")
	}

	return s.Db().Delete(&model.ArticleClass{}, "id = ? and user_id = ? and is_default = 0", cid, uid).Error
}

func (s *ArticleClassService) Sort(ctx context.Context, uid, cid, mode int) error {

	var item *model.ArticleClass
	if err := s.Db().First(&item, "id = ? and user_id = ?", cid, uid).Error; err != nil {
		return err
	}

	if mode == 1 {
		maxSort, err := s.dao.MaxSort(uid)
		if err != nil {
			return err
		}

		if maxSort == item.Sort {
			return nil
		}

		return s.Db().Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.ArticleClass{}).Where("user_id = ? and sort = ?", uid, item.Sort+1).Updates(map[string]interface{}{
				"sort": gorm.Expr("sort - 1"),
			}).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.ArticleClass{}).Where("id = ? and user_id = ?", cid, uid).Updates(map[string]interface{}{
				"sort": gorm.Expr("sort + 1"),
			}).Error; err != nil {
				return err
			}

			return nil
		})
	} else {
		minSort, err := s.dao.MinSort(uid)
		if err != nil {
			return err
		}

		if minSort == item.Sort {
			return nil
		}

		return s.Db().Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.ArticleClass{}).Where("user_id = ? and sort = ?", uid, item.Sort-1).Updates(map[string]interface{}{
				"sort": gorm.Expr("sort + 1"),
			}).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.ArticleClass{}).Where("id = ? and user_id = ?", cid, uid).Updates(map[string]interface{}{
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

	defaultClass := &model.ArticleClass{}

	err := s.Db().First(defaultClass, "id = ? and is_default = ?", uid, 1).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	s.Db().Create(&model.ArticleClass{
		UserId:    uid,
		ClassName: "默认分类",
		Sort:      1,
		IsDefault: 1,
	})
}
