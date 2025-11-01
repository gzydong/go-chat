package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/gzydong/go-chat/config"
	"github.com/gzydong/go-chat/internal/apis/handler"
	"github.com/gzydong/go-chat/internal/apis/handler/admin"
	"github.com/gzydong/go-chat/internal/apis/handler/open"
	"github.com/gzydong/go-chat/internal/apis/handler/web"
	"github.com/gzydong/go-chat/internal/apis/router"
)

type Provider struct {
	Config *config.Config
	Engine *gin.Engine
}

var ProviderSet = wire.NewSet(
	router.NewRouter,

	handler.ProviderSet, // 注入 Handler
	web.ProviderSet,     // 注入 Web Handler 依赖
	admin.ProviderSet,   // 注入 Admin Handler 依赖
	open.ProviderSet,    // 注入 Open Handler 依赖

	// Provider
	wire.Struct(new(Provider), "*"),
)
