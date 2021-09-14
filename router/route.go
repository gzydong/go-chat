package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/middleware"
)

// InitRouter 初始化配置路由
func InitRouter() *gin.Engine {
	router := gin.Default()

	// 注册跨域中间件
	router.Use(middleware.Cors())

	RegisterWebRoute(router)
	RegisterApiRoute(router)
	RegisterWsRoute(router)

	return router
}
