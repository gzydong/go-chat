package other

import (
	"github.com/urfave/cli/v2"
)

type OtherCommand *cli.Command

// Subcommands 子命令
type Subcommands struct {
}

func NewOtherCommand(subcommands *Subcommands) OtherCommand {
	return &cli.Command{
		Name:        "other",
		Usage:       "临时任务",
		Subcommands: []*cli.Command{},
	}
}
