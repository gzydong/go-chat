package main

import (
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/gateway/internal/handler"
	"go-chat/internal/gateway/internal/process"
	"go-chat/internal/provider"
)

type AppProvider struct {
	Config    *config.Config
	Engine    *gin.Engine
	Coroutine *process.Server
	Handler   *handler.Handler
	Providers *provider.Providers
}

func NewTcpServer(app *AppProvider) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", app.Config.Server.Tcp))

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		fmt.Println("RemoteAddr===>", conn.RemoteAddr())

		// TCP 分发
		go app.Handler.Dispatch(conn)
	}
}
