package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/middleware"
	"go-chat/internal/pkg/im"
	"go-chat/internal/repository/cache"

	"go-chat/config"
	"go-chat/internal/im_gateway/internal/handler"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handle *handler.Handler, session *cache.TokenSessionStorage) *gin.Engine {

	router := gin.Default()

	// 授权验证中间件
	authorize := middleware.Auth(conf.Jwt.Secret, "api", session)

	// 查看客户端连接状态
	router.GET("/wss/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, entity.H{
			"max_client_id": im.Counter.GetMaxID(),
			"chat":          im.Session.Chat.Count(),
			"example":       im.Session.Example.Count(),
		})
	})

	router.GET("/wss/default.io", authorize, ichat.HandlerFunc(handle.Chat.WsConn))
	router.GET("/wss/example.io", authorize, ichat.HandlerFunc(handle.Example.WsConnect))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, entity.H{"ok": "success"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, entity.H{"msg": "请求地址不存在"})
	})

	return router
}
