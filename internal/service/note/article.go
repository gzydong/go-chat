package note

import (
	"context"
	"errors"
	"fmt"
	"html"

	"go-chat/internal/repository/model"
	"gorm.io/gorm"

	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/service"
)

type ArticleService struct {
	*service.BaseService
}

func NewArticleService(baseService *service.BaseService) *ArticleService {
	return &ArticleService{BaseService: baseService}
}

// Detail 笔记详情
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

type ArticleEditOpt struct {
	UserId    int
	ArticleId int
	ClassId   int
	Title     string
	Content   string
	MdContent string
}

// Create 创建笔记
func (s *ArticleService) Create(ctx context.Context, opts *ArticleEditOpt) (int, error) {

	abstract := strutil.MtSubstr(opts.MdContent, 0, 200)

	abstract = strutil.Strip(abstract)

	data := &model.Article{
		UserId:   opts.UserId,
		ClassId:  opts.ClassId,
		Title:    opts.Title,
		Image:    strutil.ParseHtmlImage(opts.Content),
		Abstract: abstract,
		Status:   1,
	}

	err := s.Db().Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(data).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.ArticleDetail{
			ArticleId: data.Id,
			MdContent: html.EscapeString(opts.MdContent),
			Content:   html.EscapeString(opts.Content),
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

// Update 更新笔记信息
func (s *ArticleService) Update(ctx context.Context, opts *ArticleEditOpt) error {

	abstract := strutil.Strip(opts.MdContent)
	abstract = strutil.MtSubstr(abstract, 0, 200)

	return s.Db().Transaction(func(tx *gorm.DB) error {

		if err := tx.Model(&model.Article{}).Where("id = ? and user_id = ?", opts.ArticleId, opts.UserId).Updates(&model.Article{
			Title:    opts.Title,
			Image:    strutil.ParseHtmlImage(opts.Content),
			Abstract: abstract,
		}).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.ArticleDetail{}).Where("article_id = ?", opts.ArticleId).Updates(&model.ArticleDetail{
			MdContent: html.EscapeString(opts.MdContent),
			Content:   html.EscapeString(opts.Content),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

type ArticleListOpt struct {
	UserId   int
	Keyword  string
	FindType int
	Cid      int
	Page     int
}

// List 笔记列表
func (s *ArticleService) List(ctx context.Context, opts *ArticleListOpt) ([]*model.ArticleItem, error) {

	query := s.Db().Table("article").Select("article.*,article_class.class_name")
	query.Joins("left join article_class on article_class.id = article.class_id")
	query.Where("article.user_id = ?", opts.UserId)

	if opts.FindType == 2 {
		query.Where("article.is_asterisk = ?", 1)
	} else if opts.FindType == 3 {
		query.Where("article.class_id = ?", opts.Cid)
	} else if opts.FindType == 4 {
		query.Where("FIND_IN_SET(?,article.tags_id)", opts.Cid)
	}

	if opts.FindType == 5 {
		query.Where("article.status = ?", 2)
	} else {
		query.Where("article.status = ?", 1)
	}

	if opts.Keyword != "" {
		query.Where("article.title like ?", fmt.Sprintf("%%%s%%", opts.Keyword))
	}

	if opts.FindType == 1 {
		query.Order("article.updated_at desc").Limit(20)
	} else {
		query.Order("article.id desc")
	}

	items := make([]*model.ArticleItem, 0)
	if err := query.Scan(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

// Asterisk 笔记标记星号
func (s *ArticleService) Asterisk(ctx context.Context, uid int, articleId int, mode int) error {

	if mode != 1 {
		mode = 0
	}

	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Update("is_asterisk", mode).Error
}

// Tag 更新笔记标签
func (s *ArticleService) Tag(ctx context.Context, uid int, articleId int, tags []int) error {
	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Update("tags_id", sliceutil.IntToIds(tags)).Error
}

// Move 移动笔记分类
func (s *ArticleService) Move(ctx context.Context, uid, articleId, classId int) error {
	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Update("class_id", classId).Error
}

// UpdateStatus 修改笔记状态
func (s *ArticleService) UpdateStatus(ctx context.Context, uid int, articleId int, status int) error {
	data := map[string]interface{}{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	return s.Db().Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Updates(data).Error
}

// ForeverDelete 永久笔记
func (s *ArticleService) ForeverDelete(ctx context.Context, uid int, articleId int) error {
	var detail *model.Article
	if err := s.Db().First(&detail, "id = ? and user_id = ?", articleId, uid).Error; err != nil {
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
