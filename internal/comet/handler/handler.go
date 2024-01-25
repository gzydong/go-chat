package handler

import (
	"go-chat/config"
	"go-chat/internal/pkg/core/socket"
)

type Handler struct {
	Chat        *ChatChannel
	Example     *ExampleChannel
	Config      *config.Config
	RoomStorage *socket.RoomStorage
}
