package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/config"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type Upload struct {
	Conf *config.Config
}

func (u *Upload) Index(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	//根据当前时间鹾生成一个新的文件名
	fileNameInt := time.Now().Unix()
	fileNameStr := strconv.FormatInt(fileNameInt, 10)
	//新的文件名
	fileName := fileNameStr + path.Ext(file.Filename)
	//保存上传文件
	filePath := filepath.Join(u.Conf.Filesystem.Local.Root, "/", fileName)

	_ = c.SaveUploadedFile(file, filePath)
}
