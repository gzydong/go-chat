package v1

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/slice"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/utils"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Emoticon struct {
	filesystem *filesystem.Filesystem
}

func NewEmoticonHandler(
	filesystem *filesystem.Filesystem,
) *Emoticon {
	return &Emoticon{filesystem: filesystem}
}

// CollectList 收藏列表
func (e *Emoticon) CollectList(ctx *gin.Context) {

}

// DeleteCollect 删除收藏表情包
func (e *Emoticon) DeleteCollect(ctx *gin.Context) {

}

// Upload 上传自定义表情包
func (e *Emoticon) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("emoticon")
	if err != nil {
		response.InvalidParams(ctx, "emoticon 字段必传！")
		return
	}

	arr := []string{"png", "jpg", "jpeg", "gif"}
	ext := strings.Trim(path.Ext(file.Filename), ".")

	if !slice.InStr(ext, arr) {
		response.InvalidParams(ctx, "上传文件格式不正确,仅支持 png、jpg、jpeg 和 gif")
		return
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		response.InvalidParams(ctx, "上传文件大小不能超过5M！")
		return
	}

	open, _ := file.Open()
	defer open.Close()

	fileBytes, _ := ioutil.ReadAll(open)

	size := utils.ReadFileImage(bytes.NewReader(fileBytes))

	src := fmt.Sprintf("media/images/emoticon/%s/%s", time.Now().Format("20060102"), strutil.GenImageName(ext, size["width"], size["height"]))

	err = e.filesystem.Write(fileBytes, src)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"url": e.filesystem.PublicUrl(src),
	}, "文件上传成功")
}

// SystemList 系统表情包列表
func (e *Emoticon) SystemList(ctx *gin.Context) {

}

// SetSystemEmoticon 添加或移除系统表情包
func (e *Emoticon) SetSystemEmoticon(ctx *gin.Context) {

}
