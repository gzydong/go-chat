package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Download struct {
}

// ChatFile 下载聊天文件
func (d *Download) ChatFile() {

}

// ArticleAnnex 下载笔记附件
func (d *Download) ArticleAnnex(c *gin.Context) {
	filename := "测试中文文件名.txt"

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("测试中文文件名.txt")
}
