package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/commet"
	"go-chat/internal/job"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/logger"
)

func NewHttpCommand() ichat.Command {
	return ichat.Command{
		Name:  "http",
		Usage: "Http Command - Http API 接口服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "http")
			return apis.Run(ctx, NewHttpInjector(conf))
		},
	}
}

func NewCommetCommand() ichat.Command {
	return ichat.Command{
		Name:  "commet",
		Usage: "Commet Command - Websocket、TCP 服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "commet")
			return commet.Run(ctx, NewCommetInjector(conf))
		},
	}
}

func NewCrontabCommand() ichat.Command {
	return ichat.Command{
		Name:  "crontab",
		Usage: "Crontab Command - 定时任务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "crontab")
			return job.Cron(ctx, NewCronInjector(conf))
		},
	}
}

func NewQueueCommand() ichat.Command {
	return ichat.Command{
		Name:  "queue",
		Usage: "Queue Command - 队列任务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "queue")
			return job.RunQueue(ctx, NewQueueInjector(conf))
		},
	}
}

func NewMigrateCommand() ichat.Command {
	return ichat.Command{
		Name:  "migrate",
		Usage: "Migrate Command - 数据库初始化",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "migrate")
			return job.RunMigrate(ctx, NewMigrateInjector(conf))
		},
	}
}

func NewOtherCommand() ichat.Command {
	return ichat.Command{
		Name:  "other",
		Usage: "Other Command - 其它临时命令",
		Subcommands: []ichat.Command{
			{
				Name:  "test",
				Usage: "Test Command",
				Action: func(ctx *cli.Context, conf *config.Config) error {
					logger.InitLogger(fmt.Sprintf("%s/logs/app.log", conf.Log.Path), logger.LevelInfo, "other")
					return NewOtherInjector(conf).TestCommand.Run(ctx, conf)
				},
			},
		},
	}
}

func main() {
	app := ichat.NewApp()
	app.Register(NewHttpCommand())
	app.Register(NewCommetCommand())
	app.Register(NewCrontabCommand())
	app.Register(NewQueueCommand())
	app.Register(NewOtherCommand())
	app.Register(NewMigrateCommand())
	app.Run()
}
