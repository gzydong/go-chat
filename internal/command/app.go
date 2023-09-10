package main

import (
	"go-chat/config"
	"go-chat/internal/command/internal/command"
)

type AppProvider struct {
	Config   *config.Config
	Commands *command.Commands
}
