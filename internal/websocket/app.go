package main

import (
	"fmt"
	"net"

	"go-chat/config"
	"go-chat/internal/provider"
	"go-chat/internal/websocket/internal/handler"
	"go-chat/internal/websocket/internal/process"
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

		go app.Handler.AcceptTcp(conn)
	}
}
