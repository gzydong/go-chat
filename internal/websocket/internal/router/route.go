package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/im"
	"go-chat/internal/repository/cache"

	"go-chat/config"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/websocket/internal/handler"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handle *handler.Handler, session *cache.SessionStorage) *gin.Engine {

	router := gin.Default()

	// 授权验证中间件
	authorize := jwt.Auth(conf.Jwt.Secret, "api", session)

	// 查看客户端连接状态
	router.GET("/wss/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, entity.H{
			"server_port":   conf.App.Port,
			"max_client_id": im.Counter.GetMaxID(),
			"default":       entity.H{"online_total": im.Session.Default.Count()},
			"example":       entity.H{"online_total": im.Session.Example.Count()},
		})
	})

	router.GET("/wss/default.io", authorize, ichat.HandlerFunc(handle.DefaultWebSocket.Connect))
	router.GET("/wss/example.io", authorize, ichat.HandlerFunc(handle.ExampleWebsocket.Connect))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, entity.H{"ok": "success"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, entity.H{"msg": "请求地址不存在"})
	})

	return router
}
