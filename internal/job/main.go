package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	providers := Initialize(ctx)

	cmd := cli.NewApp()

	cmd.Name = "GoChat 脚本任务"
	cmd.Usage = "命令行管理工具"
	cmd.Commands = providers.Commands.ToCommands()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	if err := cmd.RunContext(ctx, os.Args); err != nil {
		fmt.Printf("Command Error : %s", err.Error())
	}

	time.Sleep(3 * time.Second)
}
