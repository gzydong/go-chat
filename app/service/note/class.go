package note

import (
	"context"
	"errors"
	"go-chat/app/model"
	"go-chat/app/service"
	"gorm.io/gorm"
	"time"
)

type ArticleClassService struct {
	*service.BaseService
}

func NewArticleClassService(baseService *service.BaseService) *ArticleClassService {
	return &ArticleClassService{BaseService: baseService}
}

func (s *ArticleClassService) List(ctx context.Context, uid int) {

}

func (s *ArticleClassService) Create(ctx context.Context, uid int, name string) (int, error) {

	data := &model.ArticleClass{
		UserId:    uid,
		ClassName: name,
		Sort:      1,
		CreatedAt: int(time.Now().Unix()),
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
