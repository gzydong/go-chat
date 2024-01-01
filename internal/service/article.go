package service

import (
	"context"
	"errors"
	"fmt"

	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"

	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
)

var _ IArticleService = (*ArticleService)(nil)

type IArticleService interface {
	Detail(ctx context.Context, uid int, articleId int) (*model.Article, error)
	Create(ctx context.Context, opt *ArticleEditOpt) (int, error)
	Update(ctx context.Context, opt *ArticleEditOpt) error
	List(ctx context.Context, opt *ArticleListOpt) ([]*model.ArticleListItem, error)
	Asterisk(ctx context.Context, uid int, articleId int, mode int) error
	Tag(ctx context.Context, uid int, articleId int, tags []int32) error
	Move(ctx context.Context, uid int, articleId int, classId int) error
	UpdateStatus(ctx context.Context, uid int, articleId int, status int) error
	ForeverDelete(ctx context.Context, uid int, articleId int) error
}

type ArticleService struct {
	*repo.Source
	ArticleRepo *repo.Article
}

// Detail 笔记详情
func (s *ArticleService) Detail(ctx context.Context, uid, articleId int) (*model.Article, error) {
	return s.ArticleRepo.FindByWhere(ctx, "id = ? and user_id = ?", articleId, uid)
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
func (s *ArticleService) Create(ctx context.Context, opt *ArticleEditOpt) (int, error) {

	abstract := strutil.MtSubstr(opt.MdContent, 0, 200)

	data := &model.Article{
		UserId:    opt.UserId,
		ClassId:   opt.ClassId,
		Title:     opt.Title,
		Image:     "",
		Abstract:  strutil.Strip(abstract),
		Status:    1,
		MdContent: opt.MdContent,
	}

	images := strutil.ParseMarkdownImages(opt.MdContent)
	if len(images) > 0 {
		data.Image = images[0]
	}

	err := s.Source.Db().WithContext(ctx).Create(data).Error
	if err != nil {
		return 0, err
	}

	return data.Id, nil
}

// Update 更新笔记信息
func (s *ArticleService) Update(ctx context.Context, opt *ArticleEditOpt) error {
	abstract := strutil.MtSubstr(opt.MdContent, 0, 200)

	data := map[string]any{
		"title":      opt.Title,
		"image":      "",
		"abstract":   strutil.Strip(abstract),
		"md_content": opt.MdContent,
	}

	images := strutil.ParseMarkdownImages(opt.MdContent)
	if len(images) > 0 {
		data["image"] = images[0]
	}

	_, err := s.ArticleRepo.UpdateWhere(ctx, data, "id = ? and user_id = ?", opt.ArticleId, opt.UserId)
	return err
}

type ArticleListOpt struct {
	UserId   int
	Keyword  string
	FindType int
	Cid      int
	Page     int
}

// List 笔记列表
func (s *ArticleService) List(ctx context.Context, opt *ArticleListOpt) ([]*model.ArticleListItem, error) {

	query := s.Source.Db().WithContext(ctx).Table("article").Select("article.*,article_class.class_name")
	query.Joins("left join article_class on article_class.id = article.class_id")
	query.Where("article.user_id = ?", opt.UserId)

	if opt.FindType == 2 {
		query.Where("article.is_asterisk = ?", 1)
	} else if opt.FindType == 3 {
		query.Where("article.class_id = ?", opt.Cid)
	} else if opt.FindType == 4 {
		query.Where("FIND_IN_SET(?,article.tags_id)", opt.Cid)
	}

	if opt.FindType == 5 {
		query.Where("article.status = ?", 2)
	} else {
		query.Where("article.status = ?", 1)
	}

	if opt.Keyword != "" {
		query.Where("article.title like ?", fmt.Sprintf("%%%s%%", opt.Keyword))
	}

	if opt.FindType == 1 {
		query.Order("article.updated_at desc").Limit(20)
	} else {
		query.Order("article.id desc")
	}

	items := make([]*model.ArticleListItem, 0)
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

	return s.Source.Db().WithContext(ctx).Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Update("is_asterisk", mode).Error
}

// Tag 更新笔记标签
func (s *ArticleService) Tag(ctx context.Context, uid int, articleId int, tags []int32) error {
	return s.Source.Db().WithContext(ctx).Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Update("tags_id", sliceutil.ToIds(tags)).Error
}

// Move 移动笔记分类
func (s *ArticleService) Move(ctx context.Context, uid, articleId, classId int) error {
	return s.Source.Db().WithContext(ctx).Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Update("class_id", classId).Error
}

// UpdateStatus 修改笔记状态
func (s *ArticleService) UpdateStatus(ctx context.Context, uid int, articleId int, status int) error {
	data := map[string]any{
		"status": status,
	}

	if status == 2 {
		data["deleted_at"] = timeutil.DateTime()
	}

	return s.Source.Db().WithContext(ctx).Model(&model.Article{}).Where("id = ? and user_id = ?", articleId, uid).Updates(data).Error
}

// ForeverDelete 永久笔记
func (s *ArticleService) ForeverDelete(ctx context.Context, uid int, articleId int) error {
	detail, err := s.ArticleRepo.FindByWhere(ctx, "id = ? and user_id = ?", articleId, uid)
	if err != nil {
		return err
	}

	if detail.Status != 2 {
		return errors.New("文章不能被删除")
	}

	db := s.Source.Db().WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.Article{}, detail.Id).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.ArticleAnnex{}, "user_id = ? and article_id = ?", uid, detail.Id).Error; err != nil {
			return err
		}

		return nil
	})
}
