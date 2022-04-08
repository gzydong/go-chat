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
	"go-chat/internal/provider"

	_ "go-chat/internal/pkg/validation"

	_ "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	cmd := cli.NewApp()
	cmd.Name = "Http Log"
	cmd.Usage = "GoChat 即时聊天应用"

	// 设置参数
	cmd.Flags = []cli.Flag{
		// 配置文件参数
		&cli.StringFlag{Name: "config", Aliases: []string{"c"}, Value: "./config.yaml", Usage: "配置文件路径", DefaultText: "./config.yaml"},

		// 端口号参数
		&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: 9503, Usage: "设置端口号", DefaultText: "9503"},
	}

	cmd.Action = func(tx *cli.Context) error {
		ctx, cancel := context.WithCancel(tx.Context)

		defer cancel()

		// 读取配置文件
		config := provider.ReadConfig(tx.String("config"))

		// 设置服务端口号
		config.SetPort(tx.Int("port"))

		if !config.Debug() {
			gin.SetMode(gin.ReleaseMode)

			// 配置访问日志
			f, _ := os.Create(fmt.Sprintf("%s/logs/http-access.log", config.Log.Dir))
			gin.DefaultWriter = f
		}

		providers := Initialize(ctx, config)

		eg, groupCtx := errgroup.WithContext(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

		// 启动 Http
		run(c, eg, groupCtx, cancel, providers.Server)

		return nil
	}

	_ = cmd.Run(os.Args)
}

func run(c chan os.Signal, eg *errgroup.Group, ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	// 启动 http 服务
	eg.Go(func() error {
		log.Printf("HTTP Listen %s", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Listen err: %s", err)
		}

		return nil
	})

	eg.Go(func() error {
		defer func() {
			cancel()
			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer timeCancel()
			if err := server.Shutdown(timeCtx); err != nil {
				log.Fatalf("HTTP Shutdown err: %s", err)
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
	}

	log.Fatal("HTTP Shutdown")
}
