package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	cmd := cli.App{
		Name: "GoChat 脚本任务",
		Commands: []*cli.Command{
			{
				Name:  "crontab",
				Usage: "定时任务",
				Subcommands: []*cli.Command{
					{
						Name: "clear_note",
						Action: func(context *cli.Context) error {
							fmt.Println("clear_note")
							return nil
						},
					},
					{
						Name: "clear_file",
						Action: func(context *cli.Context) error {
							fmt.Println("clear_file")
							return nil
						},
					},
				},
			},
			{
				Name:        "queue",
				Usage:       "队列任务",
				Subcommands: []*cli.Command{},
			},
			{
				Name:        "other",
				Usage:       "自定义",
				Subcommands: []*cli.Command{},
			},
		},
	}

	_ = cmd.Run(os.Args)
}
