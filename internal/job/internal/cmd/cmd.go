package cmd

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/job/internal/cmd/crontab"
	"go-chat/internal/job/internal/cmd/other"
	"go-chat/internal/job/internal/cmd/queue"
	"reflect"
)

type Commands struct {
	CrontabCommand crontab.CrontabCommand
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
