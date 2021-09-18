package router

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/controller/wss"
	"go-chat/app/middleware"
)

type WsControllerGroup struct {
	WsController *wss.WsController
}

// RegisterWsRoute 注册 Websocket 路由
func RegisterWsRoute(router *gin.Engine) {
	ControllerGroup := &WsControllerGroup{}

	router.GET("/ws/socket.io", middleware.JwtAuth("api"), ControllerGroup.WsController.SocketIo)
}
