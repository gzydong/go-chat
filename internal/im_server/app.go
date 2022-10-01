package main

import (
	"fmt"
	"net"

	"go-chat/config"
	"go-chat/internal/im_server/internal/handler"
	"go-chat/internal/im_server/internal/process"
	"go-chat/internal/provider"
)

type AppProvider struct {
	Config    *config.Config
	Server    provider.WebsocketServer
	Coroutine *process.Server
	Handler   *handler.Handler
}

func NewTcpServer(app *AppProvider) {
	listener, _ := net.Listen("tcp", "0.0.0.0:9505")

	defer func() {
		_ = listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}

		go app.Handler.ConnDispatcher(conn)
	}
}
