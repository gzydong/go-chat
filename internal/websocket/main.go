package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"go-chat/config"
	"go-chat/internal/pkg/im"
	"go-chat/internal/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	cmd := cli.NewApp()

	cmd.Name = "LumenIM 在线聊天"
	cmd.Usage = "Websocket Server"

	// 设置参数
	cmd.Flags = []cli.Flag{
		// 配置文件参数
		&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.yaml", Usage: "配置文件路径", DefaultText: "./config.yaml"},

		&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: 9504, Usage: "设置端口号", DefaultText: "9504"},
	}

	cmd.Action = func(tx *cli.Context) error {
		ctx, cancel := context.WithCancel(tx.Context)
		defer cancel()

		eg, groupCtx := errgroup.WithContext(ctx)

		// 初始化 IM 渠道配置
		im.Initialize(groupCtx, eg)

		// 读取配置文件
		conf := config.ReadConfig(tx.String("config"))

		// 设置服务端口号
		conf.SetPort(tx.Int("port"))

		// 设置日志输出
		logger.SetOutput(conf.GetLogPath(), "logger-ws")

		if !conf.Debug() {
			gin.SetMode(gin.ReleaseMode)
		}

		app := Initialize(ctx, conf)

		c := make(chan os.Signal, 1)

		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		// 延时启动守护协程
		time.AfterFunc(3*time.Second, func() {
			app.Coroutine.Start(eg, groupCtx)
		})

		log.Printf("Websocket Server ID   :%s", conf.ServerId())
		log.Printf("Websocket Listen Port :%d", conf.App.Port)
		log.Printf("Websocket Server Pid  :%d", os.Getpid())

		return start(c, eg, groupCtx, cancel, app.Server)
	}

	_ = cmd.Run(os.Args)
}

func start(c chan os.Signal, eg *errgroup.Group, ctx context.Context, cancel context.CancelFunc, server *http.Server) error {

	// 启动 http 服务
	eg.Go(func() error {

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Websocket Listen Err: %s", err)
		}

		return err
	})

	eg.Go(func() error {
		defer func() {
			cancel()

			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer timeCancel()

			if err := server.Shutdown(timeCtx); err != nil {
				log.Printf("Websocket Shutdown Err: %s \n", err)
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return nil
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("eg error: %s", err)
		return err
	}

	log.Fatal("Websocket Shutdown")

	return nil
}
