package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/extend/socket"
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

	route := router.InitRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: route,
	}

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			<-ticker.C

			fmt.Printf("当前客户端连接数 :%v\n", socket.Manager.DefaultChannel.Count)
		}
	}()

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)

	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown Error:", err)
	}

	log.Println("Server Shutdown")
}
