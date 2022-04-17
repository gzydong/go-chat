package other

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/pkg/cmdutil"
)

type Command *cli.Command

// Subcommands 注册子命令
type Subcommands struct {
	TestCommand TestCommand
}

func NewOtherCommand(subcommands *Subcommands) Command {
	return &cli.Command{
		Name:        "other",
		Usage:       "临时任务",
		Subcommands: cmdutil.ToSubCommand(subcommands),
	}
}
