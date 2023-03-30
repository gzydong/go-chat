package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/http/internal/handler/admin"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/ichat/middleware"
	"go-chat/internal/repository/cache"
)

// RegisterAdminRoute 注册 Admin 路由
func RegisterAdminRoute(secret string, router *gin.Engine, handler *admin.Handler, session *cache.JwtTokenStorage) {

	// 授权验证中间件
	authorize := middleware.Auth(secret, "admin", session)

	// v1 接口
	v1 := router.Group("/admin/v1")
	{
		index := v1.Group("/index")
		{
			index.GET("", ichat.HandlerFunc(handler.V1.Index.Index))
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", ichat.HandlerFunc(handler.V1.Auth.Login))
			auth.GET("/captcha", ichat.HandlerFunc(handler.V1.Auth.Captcha))
			auth.GET("/logout", authorize, ichat.HandlerFunc(handler.V1.Auth.Logout))
			auth.POST("/refresh", authorize, ichat.HandlerFunc(handler.V1.Auth.Refresh))
		}
	}
}
