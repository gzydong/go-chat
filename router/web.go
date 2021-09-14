package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/controller/web"
)

// ApiControllerGroup 控制器分组
type WebControllerGroup struct {
	IndexController *web.IndexController
}

// RegisterWebRoute 注册 Web 路由
func RegisterWebRoute(router *gin.Engine) {
	ControllerGroup := &WebControllerGroup{}

	router.GET("/", ControllerGroup.IndexController.Index)
}
