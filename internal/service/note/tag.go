package note

import (
	"context"
	"errors"

	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type ArticleTagService struct {
	*service.BaseService
}

func NewArticleTagService(baseService *service.BaseService) *ArticleTagService {
	return &ArticleTagService{BaseService: baseService}
}

func (s *ArticleTagService) Create(ctx context.Context, uid int, tag string) (int, error) {
	data := &model.ArticleTag{
		UserId:  uid,
		TagName: tag,
		Sort:    1,
	}

	if err := s.Db().Create(data).Error; err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (s *ArticleTagService) Update(ctx context.Context, uid int, tagId int, tag string) error {
	return s.Db().Model(&model.ArticleTag{}).Where("id = ? and user_id = ?", tagId, uid).UpdateColumn("tag_name", tag).Error
}

func (s *ArticleTagService) Delete(ctx context.Context, uid int, tagId int) error {
	var num int64
	if err := s.Db().Model(&model.Article{}).Where("user_id = ? and FIND_IN_SET(?,tags_id)", uid, tagId).Count(&num).Error; err != nil {
		return err
	}

	if num > 0 {
		return errors.New("标签已被使用不能删除")
	}

	return s.Db().Delete(&model.ArticleTag{}, "id = ? and user_id = ?", tagId, uid).Error
}

func (s *ArticleTagService) List(ctx context.Context, uid int) ([]*model.TagItem, error) {
	items := make([]*model.TagItem, 0)

	err := s.Db().Model(&model.ArticleTag{}).Select("id", "tag_name").Where("user_id = ?", uid).Scan(&items).Error
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var num int64
		if err := s.Db().Model(&model.Article{}).Where("user_id = ? and status = 1 and FIND_IN_SET(?,tags_id)", uid, item.Id).Count(&num).Error; err == nil {
			item.Count = int(num)
		}
	}

	return items, nil
}
