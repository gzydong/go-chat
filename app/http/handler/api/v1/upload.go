package v1

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/strutil"
	"go-chat/config"
	"strings"
	"time"
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

func (u *Upload) FileStream(ctx *gin.Context) {
	params := &request.UploadFileStreamRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	params.Stream = strings.Replace(params.Stream, "data:image/png;base64,", "", 1)
	params.Stream = strings.Replace(params.Stream, " ", "+", 1)

	stream, _ := base64.StdEncoding.DecodeString(params.Stream)

	object := fmt.Sprintf("public/media/image/avatar/%s/%s", time.Now().Format("20060102"), strutil.GenImageName("png", 200, 200))

	err := u.filesystem.Write(stream, object)
	if err != nil {
		response.BusinessError(ctx, "文件上传失败")
		return
	}

	response.Success(ctx, gin.H{
		"avatar": u.filesystem.PublicUrl(object),
	})
}
