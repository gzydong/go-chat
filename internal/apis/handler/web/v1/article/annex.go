package article

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"time"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Annex struct {
	ArticleAnnexRepo    *repo.ArticleAnnex
	ArticleAnnexService service.IArticleAnnexService
	Filesystem          filesystem.IFilesystem
}

// Upload 上传附件
func (c *Annex) Upload(ctx *core.Context) error {
	in := &web.ArticleAnnexUploadRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
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
		return ctx.ErrorBusiness("附件上传失败")
	}

	ext := strutil.FileSuffix(file.Filename)

	filePath := fmt.Sprintf("article-files/%s/%s", time.Now().Format("200601"), strutil.GenFileName(ext))
	if err := c.Filesystem.Write(c.Filesystem.BucketPrivateName(), filePath, stream); err != nil {
		return ctx.ErrorBusiness("附件上传失败")
	}

	data := &model.ArticleAnnex{
		UserId:       ctx.UserId(),
		ArticleId:    int(in.ArticleId),
		Drive:        entity.FileDriveMode(c.Filesystem.Driver()),
		Suffix:       ext,
		Size:         int(file.Size),
		Path:         filePath,
		OriginalName: file.Filename,
		Status:       1,
		DeletedAt: sql.NullTime{
			Valid: false,
		},
	}

	if err := c.ArticleAnnexService.Create(ctx.Ctx(), data); err != nil {
		return ctx.ErrorBusiness("附件上传失败")
	}

	return ctx.Success(&web.ArticleAnnexUploadResponse{
		AnnexId:   int32(data.Id),
		AnnexSize: int32(data.Size),
		AnnexName: data.OriginalName,
		CreatedAt: data.CreatedAt.Format(time.DateTime),
	})
}

// Delete 删除附件
func (c *Annex) Delete(ctx *core.Context) error {

	in := &web.ArticleAnnexDeleteRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleAnnexService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(in.AnnexId), 2)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAnnexDeleteResponse{})
}

// Recover 恢复附件
func (c *Annex) Recover(ctx *core.Context) error {

	in := &web.ArticleAnnexRecoverRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleAnnexService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(in.AnnexId), 1)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAnnexRecoverResponse{})
}

// RecycleList 附件回收站列表
func (c *Annex) RecycleList(ctx *core.Context) error {
	items, err := c.ArticleAnnexRepo.RecoverList(ctx.Ctx(), ctx.UserId())

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	data := make([]*web.ArticleAnnexRecoverListResponse_Item, 0)

	for _, item := range items {
		at := time.Until(item.DeletedAt.Add(time.Hour * 24 * 30))

		data = append(data, &web.ArticleAnnexRecoverListResponse_Item{
			AnnexId:      int32(item.Id),
			AnnexName:    item.OriginalName,
			ArticleId:    int32(item.ArticleId),
			ArticleTitle: item.Title,
			CreatedAt:    item.CreatedAt.Format(time.DateTime),
			DeletedAt:    item.DeletedAt.Format(time.DateTime),
			Day:          int32(math.Ceil(at.Seconds() / 86400)),
		})
	}

	return ctx.Success(&web.ArticleAnnexRecoverListResponse{
		Items: data,
		Paginate: &web.Paginate{
			Page:  1,
			Size:  10000,
			Total: int32(len(data)),
		},
	})
}

// ForeverDelete 永久删除附件
func (c *Annex) ForeverDelete(ctx *core.Context) error {

	in := &web.ArticleAnnexForeverDeleteRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleAnnexService.ForeverDelete(ctx.Ctx(), ctx.UserId(), int(in.AnnexId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAnnexForeverDeleteResponse{})
}

// Download 下载笔记附件
func (c *Annex) Download(ctx *core.Context) error {

	in := &web.ArticleAnnexDownloadRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	info, err := c.ArticleAnnexRepo.FindById(ctx.Ctx(), int(in.AnnexId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	if info.UserId != ctx.UserId() {
		return ctx.Forbidden("无权限下载")
	}

	switch info.Drive {
	case entity.FileDriveLocal:
		if c.Filesystem.Driver() != filesystem.LocalDriver {
			return ctx.ErrorBusiness("未知文件驱动类型")
		}

		filePath := c.Filesystem.(*filesystem.LocalFilesystem).Path(c.Filesystem.BucketPrivateName(), info.Path)
		ctx.Context.FileAttachment(filePath, info.OriginalName)
	case entity.FileDriveMinio:
		ctx.Context.Redirect(http.StatusFound, c.Filesystem.PrivateUrl(c.Filesystem.BucketPrivateName(), info.Path, info.OriginalName, 60*time.Second))
	default:
		return ctx.ErrorBusiness("未知文件驱动类型")
	}

	return nil
}
