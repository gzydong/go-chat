package cmd

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/cmd/crontab"
	"go-chat/internal/job/internal/cmd/other"
	"go-chat/internal/job/internal/cmd/queue"
)

type Commands struct {
	CrontabCommand crontab.CrontabCommand
	QueueCommand   queue.QueueCommand
	OtherCommand   other.OtherCommand
}

func (cmd *Commands) ToCommands() []*cli.Command {
	cmds := make([]*cli.Command, 0)

	cmds = append(cmds, cmd.CrontabCommand)
	cmds = append(cmds, cmd.QueueCommand)
	cmds = append(cmds, cmd.OtherCommand)

	return cmds
}
