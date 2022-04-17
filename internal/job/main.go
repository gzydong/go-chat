package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	providers := Initialize(ctx)

	cmd := cli.NewApp()

	cmd.Name = "LumenIM 在线聊天"
	cmd.Usage = "命令行管理工具"
	cmd.Commands = providers.Commands.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	fmt.Println()

	if err := cmd.RunContext(ctx, os.Args); err != nil {
		fmt.Printf("Command Error : %s", err.Error())
	}
}
