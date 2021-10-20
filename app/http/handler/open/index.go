package open

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
	"go-chat/app/pkg/im"
	"go-chat/app/pkg/timeutil"
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
		"date":  timeutil.DateTime(),
		"ip":    c.ClientIP(),
		"websocket": gin.H{
			"default": im.Manager.DefaultChannel.Count,
		},
	})
}
