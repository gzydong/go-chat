package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/cache"
	"go-chat/app/http/handler"
	"go-chat/app/http/middleware"
	"go-chat/config"
)

// RegisterWsRoute 注册 Websocket 路由
func RegisterWsRoute(conf *config.Config, router *gin.Engine, handler *handler.Handler, session *cache.Session) {
	// 授权验证中间件
	authorize := middleware.JwtAuth(conf, "api", session)

	router.GET("/wss/default.io", authorize, handler.DefaultWebSocket.Connect)
}
