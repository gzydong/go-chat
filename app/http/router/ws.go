package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/handler/ws"
	"go-chat/app/http/middleware"
)

type WsControllerGroup struct {
	WsController *ws.WsController
}

// RegisterWsRoute 注册 Websocket 路由
func RegisterWsRoute(router *gin.Engine) {
	ControllerGroup := new(WsControllerGroup)
	router.GET("/ws/socket.io", middleware.JwtAuth("api"), ControllerGroup.WsController.SocketIo)
	router.GET("/ws/admin.io", middleware.JwtAuth("api"), ControllerGroup.WsController.AdminIo)
}
