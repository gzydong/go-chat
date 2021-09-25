package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler"
	"go-chat/app/http/middleware"
	"go-chat/config"
)

// API 授权守卫
var ApiGuard = "api"

// RegisterApiRoute 注册 API 路由
func RegisterApiRoute(conf *config.Config, router *gin.Engine, handler *handler.Handler) {
	group := router.Group("/api/v1")
	{
		// 授权相关分组
		auth := group.Group("/auth")
		{
			auth.POST("/login", handler.Auth.Login)
			auth.POST("/register", handler.Auth.Register)
		}

		// 用户相关分组
		user := group.Group("/user").Use(middleware.JwtAuth(conf, ApiGuard))
		{
			user.GET("/detail", handler.User.Detail)
		}
	}
}
