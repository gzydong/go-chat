package httpapi

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/httpapi/handler"
	"go-chat/internal/httpapi/handler/admin"
	"go-chat/internal/httpapi/handler/open"
	"go-chat/internal/httpapi/handler/web"
	"go-chat/internal/httpapi/router"
)

type AppProvider struct {
	Config *config.Config
	Engine *gin.Engine
}

var ProviderSet = wire.NewSet(
	router.NewRouter,

	handler.ProviderSet, // 注入 Handler
	web.ProviderSet,     // 注入 Web Handler 依赖
	admin.ProviderSet,   // 注入 Admin Handler 依赖
	open.ProviderSet,    // 注入 Open Handler 依赖

	// AppProvider
	wire.Struct(new(AppProvider), "*"),
)
