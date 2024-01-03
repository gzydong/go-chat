package cron

import (
	"context"
	"time"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ClearArticle struct {
	DB         *gorm.DB
	Filesystem filesystem.IFilesystem
}

func (c *ClearArticle) Name() string {
	return "clear.article"
}

// Spec 配置定时任务规则
// 每天凌晨1点执行
func (c *ClearArticle) Spec() string {
	return "0 1 * * *"
}

func (c *ClearArticle) Enable() bool {
	return true
}

func (c *ClearArticle) Handle(ctx context.Context) error {

	c.clearAnnex()

	c.clearNote()

	return nil
}

// 删除回收站文章附件
func (c *ClearArticle) clearAnnex() {
	lastId := 0
	size := 100

	for {
		items := make([]*model.ArticleAnnex, 0)

		err := c.DB.Model(&model.ArticleAnnex{}).Where("id > ? and status = 2 and deleted_at <= ?", lastId, time.Now().AddDate(0, 0, -30)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			break
		}

		for _, item := range items {
			_ = c.Filesystem.Delete(c.Filesystem.BucketPrivateName(), item.Path)
			c.DB.Delete(&model.ArticleAnnex{}, item.Id)
		}

		if len(items) < size {
			break
		}

		lastId = items[size-1].Id
	}
}

// 删除回收站笔记
func (c *ClearArticle) clearNote() {
	lastId := 0
	size := 100

	for {
		items := make([]*model.Article, 0)

		err := c.DB.Model(&model.Article{}).Where("id > ? and status = 2 and deleted_at <= ?", lastId, time.Now().AddDate(0, 0, -30)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			break
		}

		for _, item := range items {
			subItems := make([]*model.ArticleAnnex, 0)

			if err := c.DB.Model(&model.ArticleAnnex{}).Select("path").Where("article_id = ?", item.Id).Scan(&subItems).Error; err != nil {
				continue
			}

			for _, subItem := range subItems {
				_ = c.Filesystem.Delete(c.Filesystem.BucketPrivateName(), subItem.Path)

				c.DB.Delete(&model.ArticleAnnex{}, subItem.Id)
			}

			c.DB.Delete(&model.Article{}, item.Id)
		}

		if len(items) < size {
			break
		}

		lastId = items[size-1].Id
	}
}
