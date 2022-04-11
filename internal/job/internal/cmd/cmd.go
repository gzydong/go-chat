package cmd

import (
	"reflect"

	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/cmd/cron"
	"go-chat/internal/job/internal/cmd/other"
	"go-chat/internal/job/internal/cmd/queue"
)

type Commands struct {
	CrontabCommand cron.CrontabCommand
	QueueCommand   queue.QueueCommand
	OtherCommand   other.OtherCommand
}

func (cmd *Commands) ToCommands() []*cli.Command {
	commands := make([]*cli.Command, 0)

	elem := reflect.ValueOf(cmd).Elem()
	tp := reflect.TypeOf(&cli.Command{})
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Convert(tp).Interface().(*cli.Command); ok {
			commands = append(commands, v)
		}
	}

	return commands
}
