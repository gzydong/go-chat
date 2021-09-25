package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/middleware"
)

// API 授权守卫
var ApiGuard = "api"

// ApiControllerGroup 控制器分组
type ApiControllerGroup struct {
	AuthController *v1.AuthController
	UserController *v1.UserController
}

// RegisterApiRoute 注册 API 路由
func RegisterApiRoute(router *gin.Engine) {
	ControllerGroup := new(ApiControllerGroup)

	group := router.Group("/api/v1")
	{
		// 授权相关分组
		auth := group.Group("/auth")
		{
			auth.POST("/login", ControllerGroup.AuthController.Login)
			auth.POST("/register", ControllerGroup.AuthController.Register)
		}

		// 用户相关分组
		user := group.Group("/user").Use(middleware.JwtAuth(ApiGuard))
		{
			user.GET("/detail", ControllerGroup.UserController.Detail)
		}
	}
}
