package queue

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/pkg/cmdutil"
)

type Command *cli.Command

// Subcommands 注册子命令
type Subcommands struct {
}

func NewQueueCommand(subcommands *Subcommands) Command {
	return &cli.Command{
		Name:        "queue",
		Usage:       "队列任务",
		Subcommands: cmdutil.ToSubCommand(subcommands),
	}
}
