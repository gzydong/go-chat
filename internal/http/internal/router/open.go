package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/internal/http/internal/handler/open"
	"go-chat/internal/pkg/ichat"
)

// RegisterOpenRoute 注册 Open 路由
func RegisterOpenRoute(router *gin.Engine, handler *open.Handler) {

	// v1 接口
	v1 := router.Group("/open/v1")
	{
		index := v1.Group("/index")
		{
			index.GET("", ichat.HandlerFunc(handler.V1.Index.Index))
		}
	}
}
