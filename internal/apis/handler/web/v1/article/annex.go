package article

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"time"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
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
		return ctx.ErrorBusiness("附件上传失败")
	}

	ext := strutil.FileSuffix(file.Filename)

	filePath := fmt.Sprintf("article-files/%s/%s", time.Now().Format("200601"), strutil.GenFileName(ext))
	if err := c.Filesystem.Write(c.Filesystem.BucketPrivateName(), filePath, stream); err != nil {
		return ctx.ErrorBusiness("附件上传失败")
	}

	data := &model.ArticleAnnex{
		UserId:       ctx.UserId(),
		ArticleId:    int(params.ArticleId),
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
		Id:           int32(data.Id),
		Size:         int32(data.Size),
		Path:         data.Path,
		Suffix:       data.Suffix,
		OriginalName: data.OriginalName,
	})
}

// Delete 删除附件
func (c *Annex) Delete(ctx *ichat.Context) error {

	params := &web.ArticleAnnexDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleAnnexService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(params.AnnexId), 2)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAnnexDeleteResponse{})
}

// Recover 恢复附件
func (c *Annex) Recover(ctx *ichat.Context) error {

	params := &web.ArticleAnnexRecoverRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleAnnexService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(params.AnnexId), 1)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAnnexRecoverResponse{})
}

// RecoverList 附件回收站列表
func (c *Annex) RecoverList(ctx *ichat.Context) error {

	items, err := c.ArticleAnnexRepo.RecoverList(ctx.Ctx(), ctx.UserId())

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	data := make([]*web.ArticleAnnexRecoverListResponse_Item, 0)

	for _, item := range items {
		at := time.Until(item.DeletedAt.Add(time.Hour * 24 * 30))

		data = append(data, &web.ArticleAnnexRecoverListResponse_Item{
			Id:           int32(item.Id),
			ArticleId:    int32(item.ArticleId),
			Title:        item.Title,
			OriginalName: item.OriginalName,
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
func (c *Annex) ForeverDelete(ctx *ichat.Context) error {

	params := &web.ArticleAnnexForeverDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleAnnexService.ForeverDelete(ctx.Ctx(), ctx.UserId(), int(params.AnnexId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAnnexForeverDeleteResponse{})
}

// Download 下载笔记附件
func (c *Annex) Download(ctx *ichat.Context) error {

	params := &web.ArticleAnnexDownloadRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	info, err := c.ArticleAnnexRepo.FindById(ctx.Ctx(), int(params.AnnexId))
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
