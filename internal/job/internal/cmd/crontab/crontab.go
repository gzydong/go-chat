package crontab

import "github.com/urfave/cli/v2"

type CrontabCommand *cli.Command

func NewCrontabCommand(clearTmpFileCommand ClearTmpFileCommand) CrontabCommand {
	return &cli.Command{
		Name:  "crontab",
		Usage: "定时任务",
		Subcommands: []*cli.Command{
			clearTmpFileCommand,
		},
	}
}
