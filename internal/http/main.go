package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	_ "go-chat/internal/pkg/validation"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

func main() {
	cmd := cli.NewApp()
	cmd.Name = "Http Server"
	cmd.Usage = "GoChat 即时聊天应用"

	// 设置参数
	cmd.Flags = []cli.Flag{
		&cli.IntFlag{Name: "port", Aliases: []string{"p"}, Value: 9503, Usage: "设置端口号", DefaultText: "9503"},
	}

	cmd.Action = func(tx *cli.Context) error {
		ctx, cancel := context.WithCancel(tx.Context)
		defer cancel()

		providers := Initialize(ctx)

		// 设置服务端口号
		providers.Server.Addr = fmt.Sprintf("0.0.0.0:%d", tx.Int("port"))

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
		log.Printf("HTTP listen : %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP listen : %s", err)
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
				log.Printf("HTTP Server Shutdown err: %s\n", err)
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

	log.Fatal("HTTP Server Shutdown")
}
