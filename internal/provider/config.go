package provider

import (
	"go-chat/config"
)

func ReadConfig(path string) *config.Config {
	return config.Init(path)
}

func NewConfig() *config.Config {
	return config.Init("./config.yaml")
}
