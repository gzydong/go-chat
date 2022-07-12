package command

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/cmd/internal/command/cron"
	"go-chat/internal/cmd/internal/command/other"
	"go-chat/internal/cmd/internal/command/queue"
	"go-chat/internal/pkg/cmdutil"
)

type Commands struct {
	CrontabCommand cron.Command
	QueueCommand   queue.Command
	OtherCommand   other.Command
}

func (c *Commands) SubCommands() []*cli.Command {
	return cmdutil.ToSubCommand(c)
}
