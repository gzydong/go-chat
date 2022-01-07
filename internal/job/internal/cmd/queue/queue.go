package queue

import "github.com/urfave/cli/v2"

type QueueCommand *cli.Command

func NewQueueCommand() QueueCommand {
	return &cli.Command{
		Name:        "queue",
		Usage:       "队列任务",
		Subcommands: []*cli.Command{},
	}
}
