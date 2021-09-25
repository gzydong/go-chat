package handler

import (
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
)

type Handler struct {
	Auth     *v1.Auth
	User     *v1.User
	Download *v1.Download
	Index    *open.Index
	Ws       *ws.WebSocket
}
