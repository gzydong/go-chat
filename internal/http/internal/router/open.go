package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"go-chat/internal/http/internal/handler/open"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/repository/cache"
)

// RegisterOpenRoute 注册 Open 路由
func RegisterOpenRoute(conf *config.Config, router *gin.Engine, handler *open.Handler, session *cache.Session) {

	// 授权验证中间件
	authorize := jwtutil.Auth(conf.Jwt.Secret, "open", session)

	// v1 接口
	v1 := router.Group("/open/v1", authorize)
	{
		index := v1.Group("/index")
		{
			index.GET("/", ichat.HandlerFunc(handler.V1.Index.Index))
		}
	}
}
