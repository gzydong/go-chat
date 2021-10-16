package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Download struct {
}

func NewDownloadHandler() *Download {
	return &Download{}
}

// ChatFile 下载聊天文件
func (d *Download) ChatFile() {

}

// ArticleAnnex 下载笔记附件
func (d *Download) ArticleAnnex(ctx *gin.Context) {
	filename := "测试中文文件名.txt"

	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.File("./testdata/测试中文文件名.txt")
}
