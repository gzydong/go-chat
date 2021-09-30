package v1

import (
	"fmt"
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

	file, _ := c.FormFile("file")

	_ = u.Filesystem.UploadedFile(file, fmt.Sprintf("/image/%s", file.Filename))

	response.Success(c, "")
}
