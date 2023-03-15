package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/middleware"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/repository/cache"

	"go-chat/config"
	"go-chat/internal/gateway/internal/handler"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handle *handler.Handler, session *cache.TokenSessionStorage) *gin.Engine {

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, entity.H{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

	// 授权验证中间件
	authorize := middleware.Auth(conf.Jwt.Secret, "api", session)

	// 查看客户端连接状态
	router.GET("/wss/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, entity.H{
			"max_client_id": socket.Counter.GetMaxID(),
			"chat":          socket.Session.Chat.Count(),
			"example":       socket.Session.Example.Count(),
		})
	})

	router.GET("/wss/default.io", authorize, ichat.HandlerFunc(handle.Chat.Conn))
	router.GET("/wss/example.io", authorize, ichat.HandlerFunc(handle.Example.Conn))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, entity.H{"ok": "success"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, entity.H{"msg": "请求地址不存在"})
	})

	return router
}
