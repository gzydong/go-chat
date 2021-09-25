package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler"
	"go-chat/app/http/middleware"
	"go-chat/config"
)

// InitRouter 初始化配置路由
func NewRouter(conf *config.Config, handler *handler.Handler) *gin.Engine {
	router := gin.Default()

	// 注册跨域中间件
	router.Use(middleware.Cors(conf))
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, conf.Server)
	})
	router.GET("/open", handler.Index.Index)
	RegisterApiRoute(conf, router, handler)
	RegisterWsRoute(conf, router, handler)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "404",
			"message": "请求地址不存在!",
		})
	})
	return router
}
