package provider

import (
	"go-chat/app/service"
	"net/http"
)

type Services struct {
	HttpServer   *http.Server
	SocketServer *service.SocketService
}
