package main

import (
	"go-chat/config"
	"go-chat/internal/provider"
)

type Provider struct {
	Config *config.Config
	Server provider.HttpServer
}
