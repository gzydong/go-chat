package router

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/internal/commet/handler"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/middleware"
	"go-chat/internal/pkg/ichat/socket"
	"go-chat/internal/repository/cache"

	"go-chat/config"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handle *handler.Handler, storage *cache.JwtTokenStorage) *gin.Engine {

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

	// 授权验证中间件
	authorize := middleware.Auth(conf.Jwt.Secret, "api", storage)

	// 查看客户端连接状态
	router.GET("/wss/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"chat":    socket.Session.Chat.Count(),
			"example": socket.Session.Example.Count(),
		})
	})

	router.GET("/wss/default.io", authorize, ichat.HandlerFunc(handle.Chat.Conn))
	router.GET("/wss/example.io", authorize, ichat.HandlerFunc(handle.Example.Conn))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{"ok": "success"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string]any{"msg": "请求地址不存在"})
	})

	return router
}
