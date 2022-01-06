package main

import (
	"context"
	"errors"
	"go-chat/internal/pkg/im"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 初始化 IM 渠道配置，后面将 IM 独立拆分部署，Http 服务下无需加载
	im.Initialize()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	providers := Initialize(ctx)

	eg, groupCtx := errgroup.WithContext(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// 启动守护协程
	providers.Process.Run(eg, groupCtx)

	// 启动 Http
	run(c, eg, groupCtx, cancel, providers.WsServer)
}

func run(c chan os.Signal, eg *errgroup.Group, ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	// 启动 http 服务
	eg.Go(func() error {
		log.Printf("Websocket listen : %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Websocket listen : %s", err)
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
				log.Printf("Websocket Shutdown error: %s\n", err)
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

	log.Fatal("Websocket Shutdown")
}
