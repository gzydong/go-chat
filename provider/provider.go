package provider

import (
	"go-chat/app/process"
	"go-chat/config"
	"net/http"
)

type Providers struct {
	Config     *config.Config
	HttpServer *http.Server
	Process    *process.Process
}
