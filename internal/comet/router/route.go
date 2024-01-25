package router

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"go-chat/internal/comet/handler"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/core/middleware"
	"go-chat/internal/pkg/core/socket"
	"go-chat/internal/repository/cache"

	"go-chat/config"
)

// NewRouter 初始化配置路由
func NewRouter(conf *config.Config, handle *handler.Handler, storage *cache.JwtTokenStorage) *gin.Engine {

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统错误，请重试!!!"})
	}))

	// 授权验证中间件
	authorize := middleware.Auth(conf.Jwt.Secret, "api", storage)

	// 查看客户端连接状态
	router.GET("/wss/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"chat":    socket.Session.Chat.Count(),
			"example": socket.Session.Example.Count(),
			"num":     handle.RoomStorage.GetRoomNum(),
		})
	})

	router.GET("/wss/default.io", authorize, core.HandlerFunc(handle.Chat.Conn))
	router.GET("/wss/example.io", authorize, core.HandlerFunc(handle.Example.Conn))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]any{"ok": "success"})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string]any{"msg": "请求地址不存在"})
	})

	debug := router.Group("/debug")
	{
		debug.GET("/", gin.WrapF(pprof.Index))
		debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.POST("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
		debug.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		debug.GET("/block", gin.WrapH(pprof.Handler("block")))
		debug.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		debug.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		debug.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		debug.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	return router
}
