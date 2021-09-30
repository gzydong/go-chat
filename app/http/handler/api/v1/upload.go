package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"go-chat/config"
)

type Upload struct {
	Conf       *config.Config
	Filesystem *filesystem.Filesystem
}

func (u *Upload) Index(c *gin.Context) {
	response.Success(c, "")
}
