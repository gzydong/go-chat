package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/pakg/im"
	"go-chat/app/wssocket"
	"go-chat/config"
	"go-chat/router"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	if gin.Mode() != gin.DebugMode {
		f, _ := os.Create("runtime/logs/gin.log")

		// 如果需要同时将日志写入文件和控制台
		gin.DefaultWriter = io.MultiWriter(f)
	}

	route := router.NewRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: route,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go SetServerRunId()
	go StartImServer()
	go OnlineCount()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	defer cache.CloseRedis()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}

	log.Println("Server Shutdown")
}

func SetServerRunId() {
	for {
		cache.NewServerRun().SetServerRunId(config.GetServerRunId(), time.Now().Unix())
		time.Sleep(10 * time.Second)
	}
}

func StartImServer() {
	go im.Manager.DefaultChannel.SetCallbackHandler(wssocket.NewDefaultChannelHandle()).Process()
	go im.Manager.AdminChannel.SetCallbackHandler(wssocket.NewAdminChannelHandle()).Process()
}

func OnlineCount() {
	// 调试信息
	go func() {
		for {
			time.Sleep(time.Second * 5)
			fmt.Printf("【%s】当前在线人数 : %d 人\n", im.Manager.DefaultChannel.Name, im.Manager.DefaultChannel.Count)
			fmt.Printf("【%s】当前在线人数 : %d 人\n", im.Manager.AdminChannel.Name, im.Manager.AdminChannel.Count)
		}
	}()
}
