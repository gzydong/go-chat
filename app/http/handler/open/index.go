package open

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct {
}

// Index 首页
func (i *IndexController) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "go-chat",
	})
}
