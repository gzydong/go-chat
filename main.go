package main

import (
	"context"
	"errors"
	"go-chat/app/pkg/im"
	_ "go-chat/app/pkg/validation"
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
	// 初始化 IM 渠道配置，后面将 IM 独立拆分部署，Http 服务下无需加载
	im.Initialize()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	providers := Initialize(ctx)

	eg, groupCtx := errgroup.WithContext(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// 运行协程消费者
	providers.Process.Run(eg, ctx)

	// 启动 http 服务
	eg.Go(func() error {
		log.Printf("HTTP listen :%d", providers.Config.Server.Port)
		if err := providers.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP listen: %s", err)
		}

		return nil
	})

	eg.Go(func() error {
		defer func() {
			cancel()
			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer timeCancel()
			if err := providers.HttpServer.Shutdown(timeCtx); err != nil {
				log.Printf("Http Shutdown error: %s\n", err)
			}
		}()

		select {
		case <-groupCtx.Done():
			return groupCtx.Err()
		case <-c:
			return nil
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("eg error: %s", err)
	}

	log.Println("providers Shutdown")
}
