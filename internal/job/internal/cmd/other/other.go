package other

import "github.com/urfave/cli/v2"

type OtherCommand *cli.Command

func NewOtherCommand() OtherCommand {
	return &cli.Command{
		Name:        "other",
		Usage:       "临时任务",
		Subcommands: []*cli.Command{},
	}
}
