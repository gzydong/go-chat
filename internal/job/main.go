package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	providers := Initialize(ctx)

	cmd := cli.App{
		Name:     "GoChat 脚本任务",
		Usage:    "命令行管理工具",
		Commands: providers.Commands.ToCommands(),
	}

	_ = cmd.Run(os.Args)
}
