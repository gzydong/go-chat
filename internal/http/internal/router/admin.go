package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/http/internal/handler/admin"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/repository/cache"
)

// RegisterAdminRoute 注册 Admin 路由
func RegisterAdminRoute(secret string, router *gin.Engine, handler *admin.Handler, session *cache.SessionStorage) {

	// 授权验证中间件
	authorize := jwt.Auth(secret, "admin", session)

	// v1 接口
	v1 := router.Group("/admin/v1")
	{
		index := v1.Group("/index")
		{
			index.GET("", ichat.HandlerFunc(handler.V1.Index.Index))
		}

		auth := v1.Group("/auth")
		{
			auth.GET("/login", ichat.HandlerFunc(handler.V1.Auth.Login))
			auth.GET("/logout", ichat.HandlerFunc(handler.V1.Auth.Logout))
			auth.GET("/refresh", authorize, ichat.HandlerFunc(handler.V1.Auth.Refresh))
		}

		other := v1.Group("/other", authorize)
		{
			other.GET("/test", ichat.HandlerFunc(func(ctx *ichat.Context) error {
				return nil
			}))
		}
	}
}
