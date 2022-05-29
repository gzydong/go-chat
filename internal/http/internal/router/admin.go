package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/pkg/jwtutil"
)

// RegisterAdminRoute 注册 Admin 路由
func RegisterAdminRoute(conf *config.Config, router *gin.Engine, handler *handler.AdminHandler, session *cache.Session) {
	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "admin", session)

	// v1 接口
	v1 := router.Group("/admin/v1", authorize)
	{
		common := v1.Group("/common")
		{
			common.GET("/index", func(context *gin.Context) {
				context.JSON(200, "holle word")
			})
		}
	}
}
