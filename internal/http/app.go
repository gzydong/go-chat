package main

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
)

type AppProvider struct {
	Config *config.Config
	Engine *gin.Engine
}
