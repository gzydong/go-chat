package article

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/service/note"
)

type Annex struct {
	service    *note.ArticleAnnexService
	fileSystem *filesystem.Filesystem
}

func NewAnnex(service *note.ArticleAnnexService, fileSystem *filesystem.Filesystem) *Annex {
	return &Annex{service, fileSystem}
}

// Upload 上传附件
func (c *Annex) Upload(ctx *ichat.Context) error {

	params := &web.ArticleAnnexUploadRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	file, err := ctx.Context.FormFile("annex")
	if err != nil {
		return ctx.InvalidParams("annex 字段必传！")
	}

	// 判断上传文件大小（10M）
	if file.Size > 10<<20 {
		return ctx.InvalidParams("附件大小不能超过10M！")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return ctx.BusinessError("附件上传失败")
	}

	ext := strutil.FileSuffix(file.Filename)

	filePath := fmt.Sprintf("private/files/note/%s/%s", timeutil.DateNumber(), strutil.GenFileName(ext))

	if err := c.fileSystem.Default.Write(stream, filePath); err != nil {
		return ctx.BusinessError("附件上传失败")
	}

	data := &model.ArticleAnnex{
		UserId:       ctx.UserId(),
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

	if err := c.service.Create(ctx.RequestCtx(), data); err != nil {
		return ctx.BusinessError("附件上传失败")
	}

	return ctx.Success(entity.H{
		"id":            data.Id,
		"size":          data.Size,
		"path":          data.Path,
		"original_name": data.OriginalName,
		"suffix":        data.Suffix,
	})
}

// Delete 删除附件
func (c *Annex) Delete(ctx *ichat.Context) error {

	params := &web.ArticleAnnexDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.UpdateStatus(ctx.RequestCtx(), ctx.UserId(), params.AnnexId, 2)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Recover 恢复附件
func (c *Annex) Recover(ctx *ichat.Context) error {

	params := &web.ArticleAnnexRecoverRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.UpdateStatus(ctx.RequestCtx(), ctx.UserId(), params.AnnexId, 1)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(entity.H{})
}

// RecoverList 附件回收站列表
func (c *Annex) RecoverList(ctx *ichat.Context) error {

	items, err := c.service.Dao().RecoverList(ctx.RequestCtx(), ctx.UserId())

	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	data := make([]map[string]interface{}, 0)

	for _, item := range items {

		at := time.Until(item.DeletedAt.Add(time.Hour * 24 * 30))

		data = append(data, map[string]interface{}{
			"id":            item.Id,
			"article_id":    item.ArticleId,
			"title":         item.Title,
			"original_name": item.OriginalName,
			"day":           math.Ceil(at.Seconds() / 86400),
		})
	}

	return ctx.Paginate(data, 1, 10000, len(items))
}

// ForeverDelete 永久删除附件
func (c *Annex) ForeverDelete(ctx *ichat.Context) error {

	params := &web.ArticleAnnexForeverDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.ForeverDelete(ctx.RequestCtx(), ctx.UserId(), params.AnnexId); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Download 下载笔记附件
func (c *Annex) Download(ctx *ichat.Context) error {

	params := &web.ArticleAnnexDownloadRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	info, err := c.service.Dao().FindById(ctx.RequestCtx(), params.AnnexId)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	if info.UserId != ctx.UserId() {
		return ctx.Unauthorized("无权限下载")
	}

	switch info.Drive {
	case entity.FileDriveLocal:
		ctx.Context.FileAttachment(c.fileSystem.Local.Path(info.Path), info.OriginalName)
	case entity.FileDriveCos:
		ctx.Context.Redirect(http.StatusFound, c.fileSystem.Cos.PrivateUrl(info.Path, 60))
	default:
		return ctx.BusinessError("未知文件驱动类型")
	}

	return nil
}
