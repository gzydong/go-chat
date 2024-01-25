package cron

import (
	"context"
	"path"
	"time"

	"go-chat/internal/pkg/core/crontab"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

var _ crontab.ICrontab = (*ClearTmpFile)(nil)

type ClearTmpFile struct {
	DB         *gorm.DB
	Filesystem filesystem.IFilesystem
}

// Spec 配置定时任务规则
// 每天凌晨1点10分执行
func (c *ClearTmpFile) Spec() string {
	return "0 0 * * *"
}

func (c *ClearTmpFile) Name() string {
	return "tmp.file.clear"
}

func (c *ClearTmpFile) Enable() bool {
	return true
}

func (c *ClearTmpFile) Do(ctx context.Context) error {
	lastId, size := 0, 100

	for {
		items := make([]*model.FileUpload, 0)

		err := c.DB.Model(&model.FileUpload{}).Where("id > ? and type = 1 and created_at <= ?", lastId, time.Now().AddDate(0, 0, -1)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			if item.Drive == 2 {
				_ = c.Filesystem.AbortMultipartUpload(c.Filesystem.BucketPrivateName(), item.Path, item.UploadId)

				c.DB.Delete(model.FileUpload{}, "user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId)
			} else {
				list := make([]*model.FileUpload, 0)
				c.DB.Table("file_upload").Where("user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId).Scan(&list)

				for _, value := range list {
					_ = c.Filesystem.Delete(c.Filesystem.BucketPublicName(), value.Path)
					c.DB.Delete(model.FileUpload{}, value.Id)
				}

				if len(list) > 0 {
					// 删除空文件夹
					_ = c.Filesystem.Delete(c.Filesystem.BucketPrivateName(), path.Dir(item.Path))
				}
			}

			if err := c.Filesystem.Delete(c.Filesystem.BucketPrivateName(), item.Path); err == nil {
				c.DB.Delete(model.FileUpload{}, item.Id)
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
