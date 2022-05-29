package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/cache"
	"go-chat/internal/http/internal/handler"
	"go-chat/internal/pkg/jwtutil"
)

// RegisterOpenRoute 注册 Open 路由
func RegisterOpenRoute(conf *config.Config, router *gin.Engine, handler *handler.OpenHandler, session *cache.Session) {
	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "open", session)

	// v1 接口
	v1 := router.Group("/open/v1", authorize)
	{
		common := v1.Group("/common")
		{
			common.GET("/index", func(context *gin.Context) {
				context.JSON(200, "holle word")
			})
		}
	}
}
