package handler

import (
	"go-chat/internal/apis/handler/admin"
	"go-chat/internal/apis/handler/open"
	"go-chat/internal/apis/handler/web"
)

type Handler struct {
	Api   *web.Handler   // 前端接口
	Admin *admin.Handler // 后台接口
	Open  *open.Handler  // 对外接口
}
