package article

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go-chat/internal/entity"
	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/model"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/service/note"
)

type Annex struct {
	service    *note.ArticleAnnexService
	fileSystem *filesystem.Filesystem
}

func NewAnnexHandler(service *note.ArticleAnnexService, fileSystem *filesystem.Filesystem) *Annex {
	return &Annex{service, fileSystem}
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
		fmt.Println(err.Error())
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
		UserId:       jwtutil.GetUid(ctx),
		ArticleId:    params.ArticleId,
		Drive:        entity.FileDriveMode(c.fileSystem.Driver()),
		Suffix:       ext,
		Size:         int(file.Size),
		Path:         filePath,
		OriginalName: file.Filename,
		Status:       1,
		DeletedAt: sql.NullTime{
			Valid: false,
		},
	}

	if err := c.service.Create(ctx.Request.Context(), data); err != nil {
		response.BusinessError(ctx, "附件上传失败")
		return
	}

	response.Success(ctx, entity.H{
		"id":            data.Id,
		"size":          data.Size,
		"path":          data.Path,
		"original_name": data.OriginalName,
		"suffix":        data.Suffix,
	})
}

// Delete 删除附件
func (c *Annex) Delete(ctx *gin.Context) {
	params := &request.ArticleAnnexDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.UpdateStatus(ctx.Request.Context(), jwtutil.GetUid(ctx), params.AnnexId, 2)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Recover 恢复附件
func (c *Annex) Recover(ctx *gin.Context) {
	params := &request.ArticleAnnexRecoverRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.UpdateStatus(ctx.Request.Context(), jwtutil.GetUid(ctx), params.AnnexId, 1)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, entity.H{})
}

// nolint 附件回收站列表
func (c *Annex) RecoverList(ctx *gin.Context) {
	items, err := c.service.Dao().RecoverList(ctx.Request.Context(), jwtutil.GetUid(ctx))

	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	data := make([]map[string]interface{}, 0)

	for _, item := range items {

		at := item.DeletedAt.AddDate(0, 0, 30).Sub(time.Now())

		data = append(data, map[string]interface{}{
			"id":            item.Id,
			"article_id":    item.ArticleId,
			"title":         item.Title,
			"original_name": item.OriginalName,
			"day":           math.Ceil(at.Seconds() / 86400),
		})
	}

	response.Success(ctx, entity.H{"rows": data})
}

// ForeverDelete 永久删除附件
func (c *Annex) ForeverDelete(ctx *gin.Context) {
	params := &request.ArticleAnnexForeverDeleteRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := c.service.ForeverDelete(ctx.Request.Context(), jwtutil.GetUid(ctx), params.AnnexId); err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// Download 下载笔记附件
func (c *Annex) Download(ctx *gin.Context) {
	params := &request.ArticleAnnexDownloadRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	info, err := c.service.Dao().FindById(ctx.Request.Context(), params.AnnexId)
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	if info.UserId != jwtutil.GetUid(ctx) {
		response.Unauthorized(ctx, "无权限下载")
		return
	}

	switch info.Drive {
	case entity.FileDriveLocal:
		ctx.FileAttachment(c.fileSystem.Local.Path(info.Path), info.OriginalName)
	case entity.FileDriveCos:
		ctx.Redirect(http.StatusFound, c.fileSystem.Cos.PrivateUrl(info.Path, 60))
	default:
		response.BusinessError(ctx, "未知文件驱动类型")
	}
}
