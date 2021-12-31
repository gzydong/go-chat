package open

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-chat/app/http/response"
	"go-chat/app/pkg/timeutil"
)

type Index struct {
	rds *redis.Client
}

func NewIndexHandler(rds *redis.Client) *Index {
	return &Index{
		rds: rds,
	}
}

// Index 首页
func (i *Index) Index(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"title": "go-chat",
		"date":  timeutil.DateTime(),
	})
}
