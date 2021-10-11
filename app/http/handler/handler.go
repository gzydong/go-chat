package handler

import (
	"go-chat/app/http/handler/api/v1"
	"go-chat/app/http/handler/open"
	"go-chat/app/http/handler/ws"
)

type Handler struct {
	Common      *v1.Common
	Auth        *v1.Auth
	User        *v1.User
	TalkMessage *v1.TalkMessage
	Download    *v1.Download
	Upload      *v1.Upload
	Index       *open.Index
	Ws          *ws.WebSocket
}
