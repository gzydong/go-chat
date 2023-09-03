package main

import (
	"context"
	"errors"
	"fmt"
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
		conf := config.New(tx.String("config"))

		// 设置日志输出
		logger.InitLogger(fmt.Sprintf("%s/http.log", conf.LogPath()), logger.LevelWarn)

		if !conf.Debug() {
			gin.SetMode(gin.ReleaseMode)
		}

		app := Initialize(conf)

		eg, groupCtx := errgroup.WithContext(context.Background())
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		log.Printf("HTTP Listen Port :%d", conf.Server.Http)
		log.Printf("HTTP Server Pid  :%d", os.Getpid())

		return run(c, eg, groupCtx, app)
	}

	_ = cmd.Run(os.Args)
}

func run(c chan os.Signal, eg *errgroup.Group, ctx context.Context, app *AppProvider) error {

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.Http),
		Handler: app.Engine,
	}

	// 启动 http 服务
	eg.Go(func() error {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		defer func() {
			log.Println("Shutting down server...")

			timeCtx, timeCancel := context.WithTimeout(context.Background(), 3*time.Second)
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
		log.Fatalf("HTTP Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")

	return nil
}
