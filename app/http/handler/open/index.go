package open

import (
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
)

type Index struct {
}

// Index 首页
func (i *Index) Index(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"title": "go-chat",
		"date":  time.Now().Format("2006-01-02 15:04:05"),
		"ip":    c.ClientIP(),
	})
}
