package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Download struct {
}

// ChatFileAction 下载聊天文件
func (d *Download) ChatFileAction() {

}

// ArticleAnnexAction 下载笔记附件
func (d *Download) ArticleAnnex(c *gin.Context) {
	filename := "README.md"

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("README.md")
}
