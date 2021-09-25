package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"golang.org/x/sync/errgroup"
)

func main() {

	// 第一步：初始化配置信息
	conf := config.Init("./config.yaml")

	fmt.Println(conf)

	if gin.Mode() != gin.DebugMode {
		f, _ := os.Create("runtime/logs/gin.log")
		// 如果需要同时将日志写入文件和控制台
		gin.DefaultWriter = io.MultiWriter(f)
	}

	ctx, cancel := context.WithCancel(context.Background())
	server := Initialize(ctx, conf)
	// // 获取到服务
	// serverRun := cache.NewServerRun(client)
	//
	eg, _ := errgroup.WithContext(ctx)
	//
	// group := sync.WaitGroup{}
	// group.Add(3)
	// go func() {
	// 	defer group.Done()
	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			log.Println("SetServerRunId Stop")
	// 			return
	// 		case <-time.After(10 * time.Second):
	// 			serverRun.SetServerID(ctx, conf.Server.ServerId, time.Now().Unix())
	// 		}
	// 	}
	// }()
	//
	// go im.Manager.DefaultChannel.SetCallbackHandler(websocket.NewDefaultChannelHandle()).Process(ctx, &group)
	// go im.Manager.AdminChannel.SetCallbackHandler(websocket.NewAdminChannelHandle()).Process(ctx, &group)

	// // 监听退出
	// eg.Go(func() error {
	// 	group.Wait()
	// 	return nil
	// })

	// 启动HTTP服务
	eg.Go(func() error {
		log.Printf("HTTP listen %s", ":8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
			return server.Shutdown(timeCtx)
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("eg error: %s", err)
	}

	log.Println("Server Shutdown")
}
