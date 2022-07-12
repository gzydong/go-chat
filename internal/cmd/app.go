package main

import (
	"go-chat/config"
	"go-chat/internal/cmd/internal/command"
)

type AppProvider struct {
	Config   *config.Config
	Commands *command.Commands
}
