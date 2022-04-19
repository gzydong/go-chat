package cron

import (
	"context"
	"path"
	"time"

	"go-chat/internal/model"
	"go-chat/internal/pkg/filesystem"
	"gorm.io/gorm"
)

type ClearTmpFileHandle struct {
	db         *gorm.DB
	fileSystem *filesystem.Filesystem
}

func NewClearTmpFile(db *gorm.DB, fileSystem *filesystem.Filesystem) *ClearTmpFileHandle {
	return &ClearTmpFileHandle{db: db, fileSystem: fileSystem}
}

// Spec 配置定时任务规则
// 每天凌晨1点10分执行
func (c *ClearTmpFileHandle) Spec() string {
	return "20 1 * * *"
}

func (c *ClearTmpFileHandle) Handle(ctx context.Context) error {

	lastId, size := 0, 100

	for {
		items := make([]*model.SplitUpload, 0)

		err := c.db.Model(&model.SplitUpload{}).Where("id > ? and type = 1 and drive = 1 and created_at <= ?", lastId, time.Now().Add(-24*time.Hour)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {

			list := make([]*model.SplitUpload, 0)
			c.db.Table("split_upload").Where("user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId).Scan(&list)

			for _, value := range list {
				if err := c.fileSystem.Local.Delete(value.Path); err == nil {
					c.db.Delete(model.SplitUpload{}, value.Id)
				}
			}

			if len(list) > 0 {
				_ = c.fileSystem.Local.DeleteDir(path.Dir(list[0].Path))
			}

			if err := c.fileSystem.Local.Delete(item.Path); err == nil {
				c.db.Delete(model.SplitUpload{}, item.Id)
			}
		}

		if len(items) == size {
			lastId = items[size-1].Id
		} else {
			break
		}
	}

	return nil
}
