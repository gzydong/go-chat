package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"go-chat/config"
)

type Upload struct {
	config     *config.Config
	filesystem *filesystem.Filesystem
}

func NewUploadHandler(
	config *config.Config,
	filesystem *filesystem.Filesystem,
) *Upload {
	return &Upload{
		config:     config,
		filesystem: filesystem,
	}
}

func (u *Upload) Index(ctx *gin.Context) {
	response.Success(ctx, "")
}
