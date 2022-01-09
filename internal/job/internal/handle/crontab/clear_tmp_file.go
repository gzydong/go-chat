package crontab

import (
	"go-chat/internal/dao"
	"go-chat/internal/model"
	"go-chat/internal/pkg/filesystem"
	"path"
	"time"
)

type ClearTmpFile struct {
	dao        *dao.SplitUploadDao
	fileSystem *filesystem.Filesystem
}

func NewClearTmpFile(dao *dao.SplitUploadDao, fileSystem *filesystem.Filesystem) *ClearTmpFile {
	return &ClearTmpFile{dao: dao, fileSystem: fileSystem}
}

func (c *ClearTmpFile) Handle() error {

	lastId, size := 0, 100

	for {
		items := make([]*model.SplitUpload, 0)
		err := c.dao.Db().Model(&model.SplitUpload{}).Where("id > ? and type = 1 and drive = 1 and created_at <= ?", lastId, time.Now().Add(-24*time.Hour)).Order("id asc").Limit(size).Scan(&items).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			list := make([]*model.SplitUpload, 0)
			c.dao.Db().Table("split_upload").Where("user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId).Scan(&list)

			for _, value := range list {
				if err := c.fileSystem.Local.Delete(value.Path); err == nil {
					c.dao.Db().Delete(model.SplitUpload{}, value.Id)
				}
			}

			if len(list) > 0 {
				_ = c.fileSystem.Local.DeleteDir(path.Dir(list[0].Path))
			}

			if err := c.fileSystem.Local.Delete(item.Path); err == nil {
				c.dao.Db().Delete(model.SplitUpload{}, item.Id)
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
