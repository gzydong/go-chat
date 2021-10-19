package open

import (
	"go-chat/app/pkg/im"
	"time"

	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
)

type Index struct {
}

func NewIndexHandler() *Index {
	return &Index{}
}

// Index 首页
func (i *Index) Index(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"title": "go-chat",
		"date":  time.Now().Format("2006-01-02 15:04:05"),
		"ip":    c.ClientIP(),
		"websocket": gin.H{
			"default": im.Manager.DefaultChannel.Count,
		},
	})
}
