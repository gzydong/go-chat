package provider

import (
	"go-chat/app/process"
	"net/http"
)

type Services struct {
	HttpServer *http.Server
	ServerRun  *process.ServerRun
	Subscribe  *process.WsSubscribe
}
