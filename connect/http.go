package connect

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/config"
)

func NewHttp(conf *config.Config, handler *gin.Engine) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: handler,
	}
}
