package article

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/errorx"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/filesystem"
	"github.com/gzydong/go-chat/internal/pkg/strutil"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
)

var _ web.IArticleAnnexHandler = (*Annex)(nil)

type Annex struct {
	ArticleAnnexRepo    *repo.ArticleAnnex
	ArticleAnnexService service.IArticleAnnexService
	Filesystem          filesystem.IFilesystem
}

func (a *Annex) Upload(ctx *gin.Context, _ *web.ArticleAnnexUploadRequest) (*web.ArticleAnnexUploadResponse, error) {
	in := &web.ArticleAnnexUploadRequest{}

	value := ctx.PostForm("article_id")
	if value == "" {
		return nil, errorx.New(400, "请选择文章")
	}

	id, _ := strconv.Atoi(value)
	if id <= 0 {
		return nil, errorx.New(400, "请选择文章")
	}

	in.ArticleId = int32(id)

	file, err := ctx.FormFile("annex")
	if err != nil {
		return nil, errorx.New(400, "annex 字段必传")
	}

	// 判断上传文件大小（10M）
	if file.Size > 10<<20 {
		return nil, errorx.New(400, "附件大小不能超过10M")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return nil, err
	}

	ext := strutil.FileSuffix(file.Filename)

	filePath := fmt.Sprintf("article-files/%s/%s", time.Now().Format("200601"), strutil.GenFileName(ext))
	if err := a.Filesystem.Write(a.Filesystem.BucketPrivateName(), filePath, stream); err != nil {
		return nil, err
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx.Request.Context())
	data := &model.ArticleAnnex{
		UserId:       uid,
		ArticleId:    int(in.ArticleId),
		Drive:        entity.FileDriveMode(a.Filesystem.Driver()),
		Suffix:       ext,
		Size:         int(file.Size),
		Path:         filePath,
		OriginalName: file.Filename,
		Status:       1,
		DeletedAt: sql.NullTime{
			Valid: false,
		},
	}

	if err := a.ArticleAnnexService.Create(ctx, data); err != nil {
		return nil, err
	}

	return &web.ArticleAnnexUploadResponse{
		AnnexId:   int32(data.Id),
		AnnexSize: int32(data.Size),
		AnnexName: data.OriginalName,
		CreatedAt: data.CreatedAt.Format(time.DateTime),
	}, nil
}

func (a *Annex) Delete(ctx context.Context, in *web.ArticleAnnexDeleteRequest) (*web.ArticleAnnexDeleteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	err := a.ArticleAnnexService.UpdateStatus(ctx, uid, int(in.AnnexId), 2)
	if err != nil {
		return nil, err
	}

	return &web.ArticleAnnexDeleteResponse{}, nil
}

func (a *Annex) Recover(ctx context.Context, req *web.ArticleAnnexRecoverRequest) (*web.ArticleAnnexRecoverResponse, error) {
	err := a.ArticleAnnexService.UpdateStatus(ctx, middleware.FormContextAuthId[entity.WebClaims](ctx), int(req.AnnexId), 1)
	if err != nil {
		return nil, err
	}

	return &web.ArticleAnnexRecoverResponse{}, nil
}

func (a *Annex) ForeverDelete(ctx context.Context, req *web.ArticleAnnexForeverDeleteRequest) (*web.ArticleAnnexForeverDeleteResponse, error) {
	if err := a.ArticleAnnexService.ForeverDelete(ctx, middleware.FormContextAuthId[entity.WebClaims](ctx), int(req.AnnexId)); err != nil {
		return nil, err
	}

	return &web.ArticleAnnexForeverDeleteResponse{}, nil
}

func (a *Annex) Download(ctx *gin.Context, _ *web.ArticleAnnexDownloadRequest) (*web.ArticleAnnexDownloadResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx.Request.Context())

	in := &web.ArticleAnnexDownloadRequest{}
	if err := ctx.ShouldBind(in); err != nil {
		return nil, err
	}

	annexId, err := strconv.Atoi(ctx.DefaultQuery("annex_id", "0"))
	if err != nil {
		return nil, err
	}

	info, err := a.ArticleAnnexRepo.FindById(ctx, annexId)
	if err != nil {
		return nil, err
	}

	if info.UserId != uid {
		return nil, errorx.New(403, "无权限下载")
	}

	switch info.Drive {
	case entity.FileDriveLocal:
		if a.Filesystem.Driver() != filesystem.LocalDriver {
			return nil, errorx.New(400, "未知文件驱动类型")
		}

		filePath := a.Filesystem.(*filesystem.LocalFilesystem).Path(a.Filesystem.BucketPrivateName(), info.Path)
		ctx.FileAttachment(filePath, info.OriginalName)
	case entity.FileDriveMinio:
		ctx.Redirect(http.StatusFound, a.Filesystem.PrivateUrl(a.Filesystem.BucketPrivateName(), info.Path, info.OriginalName, 60*time.Second))
	default:
		return nil, errorx.New(400, "未知文件驱动类型")
	}

	return &web.ArticleAnnexDownloadResponse{}, nil
}

func (a *Annex) RecoverList(ctx context.Context, req *web.ArticleAnnexRecoverListRequest) (*web.ArticleAnnexRecoverListResponse, error) {
	items, err := a.ArticleAnnexRepo.RecoverList(ctx, middleware.FormContextAuthId[entity.WebClaims](ctx))
	if err != nil {
		return nil, err
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

	return &web.ArticleAnnexRecoverListResponse{
		Items: data,
		Paginate: &web.Paginate{
			Page:  1,
			Size:  10000,
			Total: int32(len(data)),
		},
	}, nil
}
