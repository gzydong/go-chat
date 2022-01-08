package main

import (
	"go-chat/config"
	"go-chat/internal/provider"
	"go-chat/internal/websocket/internal/process"
)

type Providers struct {
	Config   *config.Config
	WsServer provider.WebsocketServer
	Process  *process.Process
}
