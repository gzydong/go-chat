package main

import (
	"context"
	"errors"
	"go-chat/app/pkg/im"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go-chat/app/validator"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := Initialize(ctx)

	eg, groupCtx := errgroup.WithContext(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	// 启动服务
	eg.Go(func() error {
		log.Printf("HTTP listen :%d", server.Config.Server.Port)
		if err := server.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP listen: %s", err)
		}

		return nil
	})

	// 启动服务跑socket
	eg.Go(func() error {
		im.GroupManage.DefaultChannel.Process(ctx)
		return nil
	})

	eg.Go(func() error {
		return server.ServerRun.Handle(ctx)
	})

	eg.Go(func() error {
		return server.Subscribe.Handle(ctx)
	})

	eg.Go(func() error {
		defer func() {
			cancel()
			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer timeCancel()
			if err := server.HttpServer.Shutdown(timeCtx); err != nil {
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

	log.Println("server Shutdown")
}
