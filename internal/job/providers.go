package main

import (
	"go-chat/config"
	"go-chat/internal/job/internal/command"
)

type Providers struct {
	Config   *config.Config
	Commands *command.Commands
}
