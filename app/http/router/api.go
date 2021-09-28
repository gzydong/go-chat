package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler"
	"go-chat/app/http/middleware"
	"go-chat/config"
)

// ApiGuard API 授权守卫
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
			auth.POST("/refresh", middleware.JwtAuth(conf, ApiGuard), handler.Auth.Refresh)
			auth.POST("/logout", middleware.JwtAuth(conf, ApiGuard), handler.Auth.Logout)
			auth.POST("/forget", handler.Auth.Forget)
			auth.POST("/sms-code", handler.Auth.SmsCode)
			auth.GET("/test", handler.Auth.Test)
		}

		// 用户相关分组
		user := group.Group("/user").Use(middleware.JwtAuth(conf, ApiGuard))
		{
			user.GET("/detail", handler.User.Detail)
		}

		download := group.Group("/download").Use(middleware.JwtAuth(conf, ApiGuard))
		{
			download.GET("/chat/file", handler.Download.ArticleAnnex)
		}
	}
}
