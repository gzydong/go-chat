package main

import (
	"go-chat/config"
	"go-chat/internal/websocket/internal/process"
	"net/http"
)

type Providers struct {
	Config   *config.Config
	WsServer *http.Server
	Process  *process.Process
}
