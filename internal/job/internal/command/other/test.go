package other

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type TestCommand *cli.Command

func NewTestCommand() TestCommand {
	return &cli.Command{
		Name:  "test",
		Usage: "临时任务",
		Action: func(context *cli.Context) error {
			fmt.Println("TestCommand")
			return nil
		},
	}
}
