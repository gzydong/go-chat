package provider

import "go-chat/config"

func NewConfig() *config.Config {
	return config.Init("./config.yaml")
}
