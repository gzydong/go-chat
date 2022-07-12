package other

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/cmd/internal/handle/other"
)

type ExampleCommand *cli.Command

func NewExampleCommand(job *other.ExampleHandle) ExampleCommand {
	return &cli.Command{
		Name:  "example",
		Usage: "使用案例",
		Action: func(tx *cli.Context) error {
			return job.Handle(tx.Context)
		},
	}
}
