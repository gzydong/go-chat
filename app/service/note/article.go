package note

import (
	"context"
	"errors"
	"fmt"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service"
	"gorm.io/gorm"
	"html"
)

type ArticleService struct {
	*service.BaseService
}

func NewArticleService(baseService *service.BaseService) *ArticleService {
	return &ArticleService{BaseService: baseService}
}

func (s *ArticleService) Detail(ctx context.Context, uid, articleId int) (*model.ArticleDetailInfo, error) {
	data := &model.Article{}

	if err := s.Db().First(data, "id = ? and user_id = ?", articleId, uid).Error; err != nil {
		return nil, err
	}

	detail := &model.ArticleDetail{}

	s.Db().First(detail, "article_id = ?", articleId)

	return &model.ArticleDetailInfo{
		Id:         data.Id,
		UserId:     data.UserId,
		ClassId:    data.ClassId,
		TagsId:     data.TagsId,
		Title:      data.Title,
		Abstract:   data.Abstract,
		Image:      data.Image,
		IsAsterisk: data.IsAsterisk,
		Status:     data.Status,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
		MdContent:  html.UnescapeString(detail.MdContent),
		Content:    html.UnescapeString(detail.Content),
	}, nil
}

func (s *ArticleService) Create(ctx context.Context, uid int, req *request.ArticleEditRequest) (int, error) {

	data := &model.Article{
		UserId:     uid,
		ClassId:    req.ClassId,
		Title:      req.Title,
		Image:      strutil.ParseImage(req.Content),
		Abstract:   strutil.MtSubstr(&req.Content, 0, 200),
		IsAsterisk: 0,
		Status:     1,
	}

	err := s.Db().Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(data).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.ArticleDetail{
			ArticleId: data.Id,
			MdContent: html.EscapeString(req.MdContent),
			Content:   html.EscapeString(req.Content),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

func (s *ArticleService) Update(ctx context.Context, uid int, req *request.ArticleEditRequest) error {
	return s.Db().Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&model.Article{}).Where("id = ? and user_id = ?", req.ArticleId, uid).Updates(&model.Article{
			Title:    req.Title,
			Image:    strutil.ParseImage(req.Content),
			Abstract: strutil.MtSubstr(&req.Content, 0, 200),
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.ArticleDetail{}).Where("article_id = ?", req.ArticleId).Updates(&model.ArticleDetail{
			MdContent: html.EscapeString(req.MdContent),
			Content:   html.EscapeString(req.Content),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *ArticleService) List(ctx context.Context, uid int, req *request.ArticleListRequest) ([]*model.Article, error) {

	query := s.Db().Model(&model.Article{})
	query.Joins("left join article_class on article_class.id = article.class_id")
	query.Where("article.user_id = ?", uid)

	if req.FindType == 2 {
		query.Where("article.is_asterisk = ?", 1)
	} else if req.FindType == 3 {
		query.Where("article.class_id = ?", req.Cid)
	} else if req.FindType == 4 {
		query.Where("FIND_IN_SET(?,article.tags_id)", req.Cid)
	}

	if req.FindType == 5 {
		query.Where("article.status = ?", 2)
	} else {
		query.Where("article.status = ?", 1)
	}

	if req.Keyword != "" {
		query.Where("article.title like ?", fmt.Sprintf("%%%s%%", req.Keyword))
	}

	if req.FindType == 1 {
		query.Order("article.updated_at desc").Limit(20)
	} else {
		query.Order("article.id desc")
	}

	items := make([]*model.Article, 0)
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (s *ArticleService) Asterisk(ctx context.Context, uid int, req *request.ArticleAsteriskRequest) error {

	mode := 0
	if req.Type == 1 {
		mode = 1
	}

	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", req.ArticleId, uid).Update("is_asterisk", mode).Error
}

func (s ArticleService) UpdateTag(ctx context.Context, uid int, req *request.ArticleTagsRequest) error {
	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", req.ArticleId, uid).Update("tags_id", slice.IntToIds(req.Tags)).Error
}

func (s *ArticleService) Move(ctx context.Context, uid int, req *request.ArticleMoveRequest) error {
	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", req.ArticleId, uid).Update("class_id", req.ClassId).Error
}

func (s *ArticleService) UpdateStatus(ctx context.Context, uid int, articleId int, status int) error {
	data := map[string]interface{}{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Updates(data).Error
}

func (s *ArticleService) ForeverDelete(ctx context.Context, uid int, req *request.ArticleForeverDeleteRequest) error {
	var detail *model.Article
	if err := s.Db().First(&detail, "id = ? and user_id = ?", req.ArticleId, uid).Error; err != nil {
		return err
	}

	if detail.Status != 2 {
		return errors.New("文章不能被删除")
	}

	return s.Db().Transaction(func(tx *gorm.DB) error {

		if err := tx.Delete(&model.ArticleDetail{}, "article_id = ?", detail.Id).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.Article{}, detail.Id).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.ArticleAnnex{}, "user_id = ? and article_id = ?", uid, detail.Id).Error; err != nil {
			return err
		}

		return nil
	})
}
