package main

import (
	"fmt"
	"net"

	"go-chat/config"
	"go-chat/internal/gateway/internal/handler"
	"go-chat/internal/gateway/internal/process"
	"go-chat/internal/provider"
)

type AppProvider struct {
	Config    *config.Config
	Server    provider.WebsocketServer
	Coroutine *process.Server
	Handler   *handler.Handler
}

func NewTcpServer(app *AppProvider) {
	listener, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", app.Config.Ports.Tcp))

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		// TCP 分发
		go app.Handler.Dispatch(conn)
	}
}
