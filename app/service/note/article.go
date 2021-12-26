package note

import (
	"context"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/strutil"
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

	query.Where("user_id = ?", uid)

	query.Order("id desc")

	items := make([]*model.Article, 0)
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
