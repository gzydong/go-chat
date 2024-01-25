package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/comet"
	"go-chat/internal/mission"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/logger"
	_ "go-chat/internal/pkg/server"
)

func NewHttpCommand() core.Command {
	return core.Command{
		Name:  "http",
		Usage: "Http Command - Http API 接口服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "http")
			return apis.Run(ctx, NewHttpInjector(conf))
		},
	}
}

func NewCometCommand() core.Command {
	return core.Command{
		Name:  "comet",
		Usage: "Comet Command - Websocket、TCP 服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "comet")
			return comet.Run(ctx, NewCommetInjector(conf))
		},
	}
}

func NewCrontabCommand() core.Command {
	return core.Command{
		Name:  "crontab",
		Usage: "Crontab Command - 定时任务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "crontab")
			return mission.Cron(ctx, NewCronInjector(conf))
		},
	}
}

func NewQueueCommand() core.Command {
	return core.Command{
		Name:  "queue",
		Usage: "Queue Command - 队列任务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "group",
				Usage: "分组",
				Value: "default",
			},
		},
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "queue")
			return mission.Queue(ctx, NewQueueInjector(conf))
		},
	}
}

func NewMigrateCommand() core.Command {
	return core.Command{
		Name:  "migrate",
		Usage: "Migrate Command - 数据库初始化",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.Init(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "migrate")
			return mission.Migrate(ctx, NewMigrateInjector(conf))
		},
	}
}

func NewOtherCommand() core.Command {
	return core.Command{
		Name:  "other",
		Usage: "Other Command - 其它临时命令",
		Subcommands: []core.Command{
			{
				Name:  "test",
				Usage: "Test Command",
				Action: func(ctx *cli.Context, conf *config.Config) error {
					logger.Init(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "other")
					return NewOtherInjector(conf).TestCommand.Run(ctx, conf)
				},
			},
		},
	}
}

func main() {
	app := core.NewApp()
	app.Register(NewHttpCommand())
	app.Register(NewCometCommand())
	app.Register(NewCrontabCommand())
	app.Register(NewQueueCommand())
	app.Register(NewOtherCommand())
	app.Register(NewMigrateCommand())
	app.Run()
}
