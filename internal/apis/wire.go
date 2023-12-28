package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go-chat/config"
	"go-chat/internal/apis/handler"
	"go-chat/internal/apis/handler/admin"
	"go-chat/internal/apis/handler/open"
	"go-chat/internal/apis/handler/web"
	"go-chat/internal/apis/router"
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
