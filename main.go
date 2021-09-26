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

	"go-chat/app/pkg/im"
	"go-chat/app/service"
	_ "go-chat/app/validator"
	"go-chat/app/websocket"
	"go-chat/config"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	HttpServer   *http.Server
	SocketServer *service.SocketService
}

func main() {

	// 第一步：初始化配置信息
	conf := config.Init("./config.yaml")

	ctx, cancel := context.WithCancel(context.Background())
	server := Initialize(ctx, conf)
	eg, _ := errgroup.WithContext(ctx)

	// 启动服务(设置redis)
	eg.Go(func() error {
		return server.SocketServer.Run(ctx)
	})

	// 启动服务跑socket
	eg.Go(func() error {
		im.Manager.DefaultChannel.SetCallbackHandler(websocket.NewDefaultChannelHandle()).Process(ctx)
		return nil
	})

	eg.Go(func() error {
		im.Manager.AdminChannel.SetCallbackHandler(websocket.NewAdminChannelHandle()).Process(ctx)
		return nil
	})

	// 启动HTTP服务
	eg.Go(func() error {
		log.Printf("HTTP listen :%d", conf.Server.Port)
		if err := server.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTP listen: %s", err)
		}

		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	eg.Go(func() error {
		select {
		case <-c:
			// 退出其他服务
			cancel()

			// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
			timeCtx, timeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer timeCancel()
			return server.HttpServer.Shutdown(timeCtx)
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("eg error: %s", err)
	}

	log.Println("Server Shutdown")
}
