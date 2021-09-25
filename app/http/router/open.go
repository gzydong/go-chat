package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler/open"
)

// OpenControllerGroup 控制器分组(对外接口)
type OpenControllerGroup struct {
	IndexController *open.IndexController
}

// RegisterOpenRoute 注册 Web 路由
func RegisterOpenRoute(router *gin.Engine) {
	ControllerGroup := new(OpenControllerGroup)

	router.GET("/open", ControllerGroup.IndexController.Index)
}
