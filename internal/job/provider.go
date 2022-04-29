package main

import (
	"go-chat/config"
	"go-chat/internal/job/internal/command"
)

type Provider struct {
	Config   *config.Config
	Commands *command.Commands
}
