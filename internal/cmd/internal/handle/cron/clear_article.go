package cron

import (
	"context"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

type ClearArticle struct {
	db         *gorm.DB
	fileSystem *filesystem.Filesystem
}

func NewClearArticle(db *gorm.DB, fileSystem *filesystem.Filesystem) *ClearArticle {
	return &ClearArticle{db: db, fileSystem: fileSystem}
}

// Spec 配置定时任务规则
// 每天凌晨1点执行
func (c *ClearArticle) Spec() string {
	return "0 1 * * *"
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

		err := c.db.Model(&model.ArticleAnnex{}).Where("id > ? and status = 2 and deleted_at <= ?", lastId, time.Now().AddDate(0, 0, -30)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			break
		}

		for _, item := range items {
			if item.Drive == entity.FileDriveLocal {
				_ = c.fileSystem.Local.Delete(item.Path)
			} else if item.Drive == entity.FileDriveCos {
				_ = c.fileSystem.Cos.Delete(item.Path)
			}

			c.db.Delete(&model.ArticleAnnex{}, item.Id)
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

		err := c.db.Model(&model.Article{}).Where("id > ? and status = 2 and deleted_at <= ?", lastId, time.Now().AddDate(0, 0, -30)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			break
		}

		for _, item := range items {
			subItems := make([]*model.ArticleAnnex, 0)

			if err := c.db.Model(&model.ArticleAnnex{}).Select("drive", "path").Where("article_id = ?", item.Id).Scan(&subItems).Error; err != nil {
				continue
			}

			for _, subItem := range subItems {
				if subItem.Drive == entity.FileDriveLocal {
					_ = c.fileSystem.Local.Delete(subItem.Path)
				} else if subItem.Drive == entity.FileDriveCos {
					_ = c.fileSystem.Cos.Delete(subItem.Path)
				}

				c.db.Delete(&model.ArticleAnnex{}, subItem.Id)
			}

			c.db.Delete(&model.Article{}, item.Id)
			c.db.Delete(&model.ArticleDetail{}, "article_id = ?", item.Id)
		}

		if len(items) < size {
			break
		}

		lastId = items[size-1].Id
	}
}
