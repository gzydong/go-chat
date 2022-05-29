package main

import (
	"go-chat/config"
	"go-chat/internal/provider"
)

type AppProvider struct {
	Config *config.Config
	Server provider.HttpServer
}
