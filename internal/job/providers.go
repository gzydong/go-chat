package main

import (
	"go-chat/config"
	"go-chat/internal/job/internal/cmd"
)

type Providers struct {
	Config   *config.Config
	Commands *cmd.Commands
}
