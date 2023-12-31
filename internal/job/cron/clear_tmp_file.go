package cron

import (
	"context"
	"path"
	"time"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/repository/model"
	"gorm.io/gorm"
)

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
	return "clear.tmp.file"
}

func (c *ClearTmpFile) Enable() bool {
	return true
}

func (c *ClearTmpFile) Handle(ctx context.Context) error {
	lastId, size := 0, 100

	for {
		items := make([]*model.SplitUpload, 0)

		err := c.DB.Model(&model.SplitUpload{}).Where("id > ? and type = 1 and created_at <= ?", lastId, time.Now().AddDate(0, 0, -1)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			if item.Drive == 2 {
				_ = c.Filesystem.AbortMultipartUpload(c.Filesystem.BucketPrivateName(), item.Path, item.UploadId)

				c.DB.Delete(model.SplitUpload{}, "user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId)
			} else {
				list := make([]*model.SplitUpload, 0)
				c.DB.Table("split_upload").Where("user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId).Scan(&list)

				for _, value := range list {
					_ = c.Filesystem.Delete(c.Filesystem.BucketPublicName(), value.Path)
					c.DB.Delete(model.SplitUpload{}, value.Id)
				}

				if len(list) > 0 {
					// 删除空文件夹
					_ = c.Filesystem.Delete(c.Filesystem.BucketPrivateName(), path.Dir(item.Path))
				}
			}

			if err := c.Filesystem.Delete(c.Filesystem.BucketPrivateName(), item.Path); err == nil {
				c.DB.Delete(model.SplitUpload{}, item.Id)
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
