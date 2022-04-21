package main

import (
	"go-chat/config"
	"go-chat/internal/provider"
)

type Providers struct {
	Config *config.Config
	Server provider.HttpServer
}
