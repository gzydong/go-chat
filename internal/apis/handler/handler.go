package handler

import (
	"github.com/gzydong/go-chat/internal/apis/handler/admin"
	"github.com/gzydong/go-chat/internal/apis/handler/open"
	"github.com/gzydong/go-chat/internal/apis/handler/web"
)

type Handler struct {
	Api   *web.Handler   // 前端接口
	Admin *admin.Handler // 后台接口
	Open  *open.Handler  // 对外接口
}
