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
	Providers *provider.Providers
}

func NewTcpServer(app *AppProvider) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", app.Config.Ports.Tcp))

	if err != nil {
		panic(err)
		return
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
