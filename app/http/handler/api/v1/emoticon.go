package v1

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/helper"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Emoticon struct {
	Filesystem *filesystem.Filesystem
}

// CollectList 收藏列表
func (e *Emoticon) CollectList(c *gin.Context) {

}

// DeleteCollect 删除收藏表情包
func (e *Emoticon) DeleteCollect(c *gin.Context) {

}

// Upload 上传自定义表情包
func (e *Emoticon) Upload(c *gin.Context) {
	file, err := c.FormFile("emoticon")
	if err != nil {
		response.InvalidParams(c, "emoticon 字段必传！")
		return
	}

	arr := []string{"png", "jpg", "jpeg", "gif"}
	ext := strings.Trim(path.Ext(file.Filename), ".")

	if !helper.InStr(ext, arr) {
		response.InvalidParams(c, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(c, "上传文件大小不能超过5M！")
		return
	}

	open, _ := file.Open()
	defer open.Close()

	fileBytes, _ := ioutil.ReadAll(open)

	size := helper.ReadFileImage(bytes.NewReader(fileBytes))
	src := fmt.Sprintf("media/images/emoticon/%s/%s", time.Now().Format("20060102"), helper.GenImageName(ext, size["width"], size["height"]))

	err = e.Filesystem.Write(fileBytes, src)
	if err != nil {
		response.BusinessError(c, err)
		return
	}

	response.Success(c, gin.H{
		"url": e.Filesystem.PublicUrl(src),
	}, "文件上传成功")
}

// SystemList 系统表情包列表
func (e *Emoticon) SystemList(c *gin.Context) {

}

// SetSystemEmoticon 添加或移除系统表情包
func (e *Emoticon) SetSystemEmoticon(c *gin.Context) {

}
