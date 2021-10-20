package provider

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go-chat/app/service"
	"gorm.io/gorm"
	"net/http"
)

type Services struct {
	Db           *gorm.DB
	Rds          *redis.Client
	Logrus       *logrus.Logger
	HttpServer   *http.Server
	SocketServer *service.SocketService
}
