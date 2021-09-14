package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type IndexController struct {
}

// Index 首页
func (i *IndexController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "go-chat",
	})
}
