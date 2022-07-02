package main

import (
	"go-chat/config"
	"go-chat/internal/provider"
	"go-chat/internal/websocket/internal/process"
)

type Provider struct {
	Config    *config.Config
	WsServer  provider.WebsocketServer
	Coroutine *process.Coroutine
}
