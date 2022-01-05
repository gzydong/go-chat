package main

import (
	"go-chat/app/websocket/internal/process"
	"go-chat/config"
	"net/http"
)

type Providers struct {
	Config   *config.Config
	WsServer *http.Server
	Process  *process.Process
}
