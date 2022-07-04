package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/http/internal/handler/admin"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/repository/cache"
)

// RegisterAdminRoute 注册 Admin 路由
func RegisterAdminRoute(conf *config.Config, router *gin.Engine, handler *admin.Handler, session *cache.Session) {

	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "admin", session)

	// v1 接口
	v1 := router.Group("/admin/v1", authorize)
	{
		index := v1.Group("/index")
		{
			index.GET("/", ichat.HandlerFunc(handler.V1.Index.Index))
		}
	}
}