package main

import (
	"net/http"

	"go-chat/config"
)

type Providers struct {
	Config *config.Config
	Server *http.Server
}
