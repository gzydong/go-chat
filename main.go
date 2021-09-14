package main

import (
	"github.com/gin-gonic/gin"
	"go-chat/router"
	"io"
	"os"
)

func main() {
	if gin.Mode() != gin.DebugMode {
		f, _ := os.Create("gin.log")

		// 如果需要同时将日志写入文件和控制台
		gin.DefaultWriter = io.MultiWriter(f)
	}

	route := router.InitRouter()

	// Listen and Server in 0.0.0.0:8080
	_ = route.Run(":8080")
}
