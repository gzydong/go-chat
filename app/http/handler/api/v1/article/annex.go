package article

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/service/note"
)

type Annex struct {
	service    *note.ArticleAnnexService
	fileSystem *filesystem.Filesystem
}

func NewAnnexHandler(service *note.ArticleAnnexService, fileSystem *filesystem.Filesystem) *Annex {
	return &Annex{service, fileSystem}
}

// List 附件列表
func (c *Annex) List(ctx *gin.Context) {

}

// Upload 上传附件
func (c *Annex) Upload(ctx *gin.Context) {
	params := &request.ArticleAnnexUploadRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	file, err := ctx.FormFile("annex")
	if err != nil {
		response.InvalidParams(ctx, "annex 字段必传！")
		return
	}

	// 判断上传文件大小（10M）
	if file.Size > 10<<20 {
		response.InvalidParams(ctx, "附件大小不能超过10M！")
		return
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		response.BusinessError(ctx, "附件上传失败")
		return
	}

	ext := strutil.FileSuffix(file.Filename)

	filePath := fmt.Sprintf("private/files/note/%s/%s", timeutil.DateNumber(), strutil.GenFileName(ext))

	if err := c.fileSystem.Default.Write(stream, filePath); err != nil {
		response.BusinessError(ctx, "附件上传失败")
		return
	}

	data := &model.ArticleAnnex{
		UserId:       auth.GetAuthUserID(ctx),
		ArticleId:    params.ArticleId,
		FileSuffix:   ext,
		FileSize:     int(file.Size),
		SaveDir:      filePath,
		OriginalName: file.Filename,
		Status:       1,
	}

	if err := c.service.Create(ctx.Request.Context(), data); err != nil {
		response.BusinessError(ctx, "附件上传失败")
		return
	}

	response.Success(ctx, gin.H{
		"id":            data.Id,
		"file_size":     data.FileSize,
		"save_dir":      data.SaveDir,
		"original_name": data.OriginalName,
		"file_suffix":   data.FileSuffix,
	})
}

// Delete 删除附件
func (c *Annex) Delete(ctx *gin.Context) {

}

// Recover 恢复附件
func (c *Annex) Recover(ctx *gin.Context) {

}

// RecoverList 附件回收站列表
func (c *Annex) RecoverList(ctx *gin.Context) {

}

// ForeverDelete 永久删除附件
func (c *Annex) ForeverDelete(ctx *gin.Context) {

}
