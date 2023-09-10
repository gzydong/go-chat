package other

import (
	"github.com/urfave/cli/v2"
	"go-chat/internal/pkg/cmdutil"
)

type Command *cli.Command

// Subcommands 注册子命令
type Subcommands struct {
	ExampleCommand ExampleCommand
	MigrateCommand MigrateCommand
}

func NewOtherCommand(subcommands *Subcommands) Command {
	return &cli.Command{
		Name:        "other",
		Usage:       "Other Commands",
		Subcommands: cmdutil.ToSubCommand(subcommands),
	}
}
