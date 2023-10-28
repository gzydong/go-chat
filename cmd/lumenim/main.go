package main

import (
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/commet"
	"go-chat/internal/httpapi"
	"go-chat/internal/job"
	"go-chat/internal/pkg/logger"
)

func NewHttpCommand() Command {
	return Command{
		Name:  "http",
		Usage: "Http Command - Http API 接口服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger("./app.log", logger.LevelInfo, "http")
			return httpapi.Run(ctx, NewHttpInjector(conf))
		},
	}
}

func NewCommetCommand() Command {
	return Command{
		Name:  "commet",
		Usage: "Commet Command - Websocket、TCP 服务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger("./app.log", logger.LevelInfo, "commet")
			return commet.Run(ctx, NewCommetInjector(conf))
		},
	}
}

func NewCrontabCommand() Command {
	return Command{
		Name:  "crontab",
		Usage: "Crontab Command - 定时任务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger("./app.log", logger.LevelInfo, "crontab")
			return job.Cron(ctx, NewCronInjector(conf))
		},
	}
}

func NewQueueCommand() Command {
	return Command{
		Name:  "queue",
		Usage: "Queue Command - 队列任务",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger("./app.log", logger.LevelInfo, "queue")
			return job.RunQueue(ctx, NewQueueInjector(conf))
		},
	}
}

func NewMigrateCommand() Command {
	return Command{
		Name:  "migrate",
		Usage: "Migrate Command - 数据库初始化",
		Action: func(ctx *cli.Context, conf *config.Config) error {
			logger.InitLogger("./app.log", logger.LevelInfo, "migrate")
			return job.RunMigrate(ctx, NewMigrateInjector(conf))
		},
	}
}

func NewOtherCommand() Command {
	return Command{
		Name:  "other",
		Usage: "Other Command - 其它临时命令",
		Subcommands: []Command{
			{
				Name:  "test",
				Usage: "Test Command",
				Action: func(ctx *cli.Context, conf *config.Config) error {
					logger.InitLogger("./app.log", logger.LevelInfo, "other")
					return NewOtherInjector(conf).TestCommand.Run(ctx, conf)
				},
			},
		},
	}
}

func main() {
	app := NewApp()
	app.Register(NewHttpCommand())
	app.Register(NewCommetCommand())
	app.Register(NewCrontabCommand())
	app.Register(NewQueueCommand())
	app.Register(NewOtherCommand())
	app.Register(NewMigrateCommand())
	app.Run()
}
