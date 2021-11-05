package provider

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
)

func NewHttpServer(conf *config.Config, handler *gin.Engine) *http.Server {

	// 默认处理端口号
	if conf.Server.Port == 0 {
		conf.Server.Port = 8080
	}

	return &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", conf.Server.Port),
		Handler: handler,
	}
}
