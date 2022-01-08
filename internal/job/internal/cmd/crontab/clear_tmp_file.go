package crontab

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go-chat/internal/dao"
	"go-chat/internal/model"
	"go-chat/internal/pkg/filesystem"
	"path"
	"time"
)

type ClearTmpFileCommand *cli.Command

// nolint
func NewClearTmpFileCommand(dao *dao.SplitUploadDao, fileSystem *filesystem.Filesystem) ClearTmpFileCommand {
	return &cli.Command{
		Name:  "clear_tmp_file",
		Usage: "清除拆分上传临时文件 (注释:删除24小时前数据)",
		Action: func(context *cli.Context) error {
			lastId, size := 0, 100

			fmt.Println("正在删除拆分上传临时文件...")

			for {
				items := make([]*model.SplitUpload, 0)
				err := dao.Db().Table("split_upload").Where("id > ? and type = 1 and drive = 1 and created_at <= ?", lastId, time.Now().Add(-24*time.Hour)).Order("id asc").Limit(size).Scan(&items).Error
				if err != nil {
					return err
				}

				for _, item := range items {
					list := make([]*model.SplitUpload, 0)
					dao.Db().Table("split_upload").Where("user_id = ? and upload_id = ? and type = 2", item.UserId, item.UploadId).Scan(&list)

					for _, value := range list {
						if err := fileSystem.Local.Delete(value.Path); err == nil {
							dao.Db().Delete(model.SplitUpload{}, value.Id)
						}
					}

					if len(list) > 0 {
						_ = fileSystem.Local.DeleteDir(path.Dir(list[0].Path))
					}

					if err := fileSystem.Local.Delete(item.Path); err == nil {
						dao.Db().Delete(model.SplitUpload{}, item.Id)
					}
				}

				if len(items) == size {
					lastId = items[size-1].Id
				} else {
					break
				}
			}

			fmt.Println("删除拆分上传临时文件已完成")

			return nil
		},
	}
}
