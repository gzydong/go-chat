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
	"go-chat/internal/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func main() {
	cmd := cli.NewApp()

	cmd.Name = "LumenIM 在线聊天"
	cmd.Usage = "Http Server"

	// 设置参数
	cmd.Flags = []cli.Flag{
		// 配置文件参数
		&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.yaml", Usage: "配置文件路径", DefaultText: "./config.yaml"},
	}

	cmd.Action = func(tx *cli.Context) error {
		// 读取配置文件
		conf := config.ReadConfig(tx.String("config"))

		// 设置日志输出
		logger.SetOutput(conf.GetLogPath(), "logger-http")

		if !conf.Debug() {
			gin.SetMode(gin.ReleaseMode)
		}

		app := Initialize(conf)

		eg, groupCtx := errgroup.WithContext(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		log.Printf("HTTP Listen Port :%d", conf.Ports.Http)
		log.Printf("HTTP Server Pid  :%d", os.Getpid())

		return run(c, eg, groupCtx, app.Server)
	}

	_ = cmd.Run(os.Args)
}

func run(c chan os.Signal, eg *errgroup.Group, ctx context.Context, server *http.Server) error {
	// 启动 http 服务
	eg.Go(func() error {

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server Listen Err: %s", err)
		}

		return err
	})

	eg.Go(func() error {
		defer func() {
			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer timeCancel()

			if err := server.Shutdown(timeCtx); err != nil {
				log.Fatalf("HTTP Server Shutdown Err: %s", err)
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
		log.Fatalf("Error: %s", err)
	}

	log.Fatal("HTTP Server Shutdown")

	return nil
}
