package crontab

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

type ClearTmpFileCommand *cli.Command

func NewClearTmpFileCommand() ClearTmpFileCommand {
	return &cli.Command{
		Name:  "clear_tmp_file",
		Usage: "清除拆分上传临时文件",
		Action: func(context *cli.Context) error {
			fmt.Println("clear_tmp_file")
			return nil
		},
	}
}
