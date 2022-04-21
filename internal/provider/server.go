package provider

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
)

type HttpServer *http.Server

type WebsocketServer *http.Server

func NewHttpServer(conf *config.Config, handler *gin.Engine) HttpServer {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.App.Port),
		Handler: handler,
	}
}

func NewWebsocketServer(conf *config.Config, handler *gin.Engine) WebsocketServer {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.App.Port),
		Handler: handler,
	}
}
