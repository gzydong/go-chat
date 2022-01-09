package crontab

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/handle/crontab"
)

type ClearArticleCommand *cli.Command

func NewClearArticleCommand(job *crontab.ClearArticle) ClearArticleCommand {
	return &cli.Command{
		Name:  "clear_article",
		Usage: "清除回收站中的笔记 (注释:删除30天前回收站笔记)",
		Action: func(context *cli.Context) error {
			var err error

			fmt.Println("[已开始] 清除回收站中的笔记")

			err = job.Handle()

			fmt.Println("[已完成] 清除回收站中的笔记")

			return err
		},
	}
}
