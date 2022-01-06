package open

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/timeutil"
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
