package main

import (
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/apis"
	"go-chat/internal/mission"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/longnet"
	_ "go-chat/internal/pkg/server"
)

// Version 服务版本号（默认）
// 构建时传入版本号
// go build -o lumenim -ldflags "-X main.Version=${IMAGE_TAG}" ./cmd/lumenim
var Version = "1.0.0"

func main() {
	app := core.NewApp(Version)
	app.Register(NewHttpCommand)
	app.Register(NewCometCommand)
	app.Register(NewCrontabCommand)
	app.Register(NewQueueCommand)
	app.Register(NewTempCommand)
	app.Register(NewMigrateCommand)
	app.Run()
}

func NewHttpCommand() core.Command {
	return core.Command{
		Name:  "http",
		Usage: "HttpAddr Command - HttpAddr API 接口服务",
		Action: func(ctx *cli.Context, c *config.Config) error {
			logger.Init(c.Log.LogFilePath("app.log"), logger.LevelInfo, "http")
			return apis.NewServer(ctx, NewHttpInjector(c))
		},
	}
}

func NewCometCommand() core.Command {
	return core.Command{
		Name:  "comet",
		Usage: "Comet Command - WebsocketAddr、TCP 服务",
		Action: func(ctx *cli.Context, c *config.Config) error {
			logger.Init(c.Log.LogFilePath("app.log"), logger.LevelInfo, "comet")
			injector := NewCometInjector(c)
			return injector.Server.Start(ctx.Context)
		},
	}
}

func NewCrontabCommand() core.Command {
	return core.Command{
		Name:  "crontab",
		Usage: "Crontab Command - 定时任务",
		Action: func(ctx *cli.Context, c *config.Config) error {
			logger.Init(c.Log.LogFilePath("app.log"), logger.LevelInfo, "crontab")
			return mission.Cron(ctx, NewCronInjector(c))
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
		Action: func(ctx *cli.Context, c *config.Config) error {
			logger.Init(c.Log.LogFilePath("app.log"), logger.LevelInfo, "queue")
			return mission.Queue(ctx, NewQueueInjector(c))
		},
	}
}

func NewMigrateCommand() core.Command {
	return core.Command{
		Name:  "migrate",
		Usage: "Migrate Command - 数据库初始化",
		Action: func(ctx *cli.Context, c *config.Config) error {
			logger.Init(c.Log.LogFilePath("app.log"), logger.LevelInfo, "migrate")
			return mission.Migrate(ctx, NewMigrateInjector(c))
		},
	}
}

func NewTempCommand() core.Command {
	return core.Command{
		Name:  "temp",
		Usage: "Temp Command - 临时命令",
		Subcommands: []core.Command{
			{
				Name:  "socket",
				Usage: "Test Command",
				Action: func(ctx *cli.Context, c *config.Config) error {
					serv := longnet.New(longnet.Options{
						WSSConfig: &longnet.WSSConfig{
							Addr: ":9501",
							Path: "/wss",
						},
					})

					serv.SetHandler(longnet.NewHandler(
						longnet.WithOpenHandler(func(smg longnet.ISessionManager, c longnet.ISession) {

						}),
						longnet.WithCloseHandler(func(cid int64, uid int64) {

						}),
						longnet.WithMessageHandler(func(smg longnet.ISessionManager, c longnet.ISession, data []byte) {
							_ = c.Write(data)
						})))

					return serv.Start(ctx.Context)
				},
			},
		},
	}
}
