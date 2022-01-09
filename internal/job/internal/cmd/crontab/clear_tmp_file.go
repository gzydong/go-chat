package crontab

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/handle/crontab"
)

type ClearTmpFileCommand *cli.Command

// nolint
func NewClearTmpFileCommand(job *crontab.ClearTmpFile) ClearTmpFileCommand {
	return &cli.Command{
		Name:  "clear_tmp_file",
		Usage: "清除拆分上传临时文件 (注释:删除24小时前数据)",
		Action: func(context *cli.Context) error {
			var err error

			fmt.Println("[已开始] 清除拆分上传临时文件")

			err = job.Handle()

			fmt.Println("[已完成] 清除拆分上传临时文件")

			return err
		},
	}
}
