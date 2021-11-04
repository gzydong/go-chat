package provider

import (
	"go-chat/app/process"
	"go-chat/config"
	"net/http"
)

type Services struct {
	Config     *config.Config
	HttpServer *http.Server
	ServerRun  *process.ServerRun
	Subscribe  *process.WsSubscribe
}
