package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/websocket/internal/handler"
	"go-chat/internal/websocket/internal/middleware"
	"net/http"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handle *handler.Handler, session *cache.Session) *gin.Engine {

	router := gin.Default()

	// 授权验证中间件
	authorize := middleware.JwtAuth(conf.Jwt.Secret, "api", session)

	router.GET("/wss/default.io", authorize, handle.DefaultWebSocket.Connect)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": "success"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"msg": "请求地址不存在"})
	})

	return router
}
