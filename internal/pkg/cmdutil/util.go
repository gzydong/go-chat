package cmdutil

import (
	"reflect"

	"github.com/urfave/cli/v2"
)

func ToSubCommand(value any) []*cli.Command {

	commands := make([]*cli.Command, 0)

	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return commands
	}

	elem := reflect.ValueOf(value).Elem()
	of := reflect.TypeOf(&cli.Command{})
	for i := 0; i < elem.NumField(); i++ {
		if v, ok := elem.Field(i).Convert(of).Interface().(*cli.Command); ok {
			commands = append(commands, v)
		}
	}

	return commands
}
