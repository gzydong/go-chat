package open

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Index struct {
}

// Index 首页
func (i *Index) Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"title": "go-chat",
	})
}
