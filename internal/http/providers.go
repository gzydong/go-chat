package main

import (
	"go-chat/config"
	"net/http"
)

type Providers struct {
	Config     *config.Config
	HttpServer *http.Server
}
