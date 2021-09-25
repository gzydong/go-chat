package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/middleware"
	"net/http"
)

// InitRouter 初始化配置路由
func NewRouter() *gin.Engine {
	router := gin.Default()

	// 注册跨域中间件
	router.Use(middleware.Cors())

	defaultRouter(router)

	RegisterApiRoute(router)
	RegisterWsRoute(router)
	RegisterOpenRoute(router)

	return router
}

func defaultRouter(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"AppName": "Lumen IM(golang)",
			"Version": "v1.0.0",
			"Author":  "837215079@qq.com",
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "404",
			"message": "请求地址不存在!",
		})
	})
}
