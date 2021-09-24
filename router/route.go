package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/middleware"
	"net/http"
)

// InitRouter 初始化配置路由
func NewRouter() *gin.Engine {
	router := gin.Default()

	// 注册跨域中间件
	router.Use(middleware.Cors())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"AppName": "Lumen IM(golang)",
			"Version": "v1.0.0",
			"Author":  "837215079@qq.com",
		})
	})

	RegisterApiRoute(router)
	RegisterWsRoute(router)
	RegisterOpenRoute(router)

	return router
}
