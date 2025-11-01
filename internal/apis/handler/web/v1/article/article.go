package article

import (
	"context"
	"html"
	"math"
	"time"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/filesystem"
	"github.com/gzydong/go-chat/internal/pkg/sliceutil"
	"github.com/gzydong/go-chat/internal/pkg/timeutil"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

var _ web.IArticleHandler = (*Article)(nil)

type Article struct {
	Source              *repo.Source
	ArticleAnnexRepo    *repo.ArticleAnnex
	ArticleClassRepo    *repo.ArticleClass
	ArticleRepo         *repo.Article
	ArticleService      service.IArticleService
	ArticleAnnexService service.IArticleAnnexService
	Filesystem          filesystem.IFilesystem
}

func (a *Article) Edit(ctx context.Context, in *web.ArticleEditRequest) (*web.ArticleEditResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	uid := session.GetAuthID()

	opt := &service.ArticleEditOpt{
		UserId:    uid,
		ArticleId: int(in.ArticleId),
		ClassId:   int(in.ClassifyId),
		Title:     in.Title,
		MdContent: in.MdContent,
	}

	if in.ArticleId == 0 {
		id, err := a.ArticleService.Create(ctx, opt)
		if err == nil {
			in.ArticleId = int32(id)
		}
	} else {
		err := a.ArticleService.Update(ctx, opt)
		if err != nil {
			return nil, err
		}
	}

	var info *model.Article
	if err := a.Source.Db().First(&info, in.ArticleId).Error; err != nil {
		return nil, err
	}

	return &web.ArticleEditResponse{
		ArticleId: int32(info.Id),
		Title:     info.Title,
		Abstract:  info.Abstract,
		Image:     info.Image,
	}, nil
}

func (a *Article) Detail(ctx context.Context, in *web.ArticleDetailRequest) (*web.ArticleDetailResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	uid := session.GetAuthID()

	detail, err := a.ArticleService.Detail(ctx, uid, int(in.ArticleId))
	if err != nil {
		return nil, err
	}

	tags := make([]*web.ArticleDetailResponse_Tag, 0)
	for _, id := range sliceutil.ParseIds(detail.TagsId) {
		tags = append(tags, &web.ArticleDetailResponse_Tag{Id: int32(id)})
	}

	files := make([]*web.ArticleDetailResponse_AnnexFile, 0)
	items, err := a.ArticleAnnexRepo.AnnexList(ctx, uid, int(in.ArticleId))
	if err == nil {
		for _, item := range items {
			files = append(files, &web.ArticleDetailResponse_AnnexFile{
				AnnexId:   int32(item.Id),
				AnnexName: item.OriginalName,
				AnnexSize: int32(item.Size),
				CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
			})
		}
	}

	return &web.ArticleDetailResponse{
		ArticleId:  int32(detail.Id),
		ClassifyId: int32(detail.ClassId),
		Title:      detail.Title,
		MdContent:  html.UnescapeString(detail.MdContent),
		IsAsterisk: int32(detail.IsAsterisk),
		CreatedAt:  timeutil.FormatDatetime(detail.CreatedAt),
		UpdatedAt:  timeutil.FormatDatetime(detail.UpdatedAt),
		TagIds:     tags,
		AnnexList:  files,
	}, nil
}

func (a *Article) List(ctx context.Context, in *web.ArticleListRequest) (*web.ArticleListResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)

	uid := session.GetAuthID()

	items, err := a.ArticleService.List(ctx, &service.ArticleListOpt{
		UserId:     uid,
		FindType:   int(in.FindType),
		Keyword:    in.Keyword,
		ClassifyId: int(in.ClassifyId),
		TagId:      int(in.TagId),
	})
	if err != nil {
		return nil, err
	}

	list := make([]*web.ArticleListResponse_Item, 0)
	for _, item := range items {
		list = append(list, &web.ArticleListResponse_Item{
			ArticleId:  int32(item.Id),
			ClassifyId: int32(item.ClassId),
			TagsId:     item.TagsId,
			Title:      item.Title,
			ClassName:  item.ClassName,
			Image:      item.Image,
			IsAsterisk: int32(item.IsAsterisk),
			Status:     int32(item.Status),
			CreatedAt:  timeutil.FormatDatetime(item.CreatedAt),
			UpdatedAt:  timeutil.FormatDatetime(item.UpdatedAt),
			Abstract:   item.Abstract,
		})
	}

	return &web.ArticleListResponse{
		Items: list,
	}, nil
}

func (a *Article) Delete(ctx context.Context, in *web.ArticleDeleteRequest) (*web.ArticleDeleteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	err := a.ArticleService.UpdateStatus(ctx, uid, int(in.ArticleId), 2)
	if err != nil {
		return nil, err
	}

	return &web.ArticleDeleteResponse{}, nil
}

func (a *Article) Recover(ctx context.Context, in *web.ArticleRecoverRequest) (*web.ArticleRecoverResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	err := a.ArticleService.UpdateStatus(ctx, uid, int(in.ArticleId), 1)
	if err != nil {
		return nil, err
	}

	return &web.ArticleRecoverResponse{}, nil
}

func (a *Article) ForeverDelete(ctx context.Context, in *web.ArticleForeverDeleteRequest) (*web.ArticleForeverDeleteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := a.ArticleService.ForeverDelete(ctx, uid, int(in.ArticleId)); err != nil {
		return nil, err
	}

	return &web.ArticleForeverDeleteResponse{}, nil
}

func (a *Article) Move(ctx context.Context, in *web.ArticleMoveRequest) (*web.ArticleMoveResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if err := a.ArticleService.Move(ctx, uid, int(in.ArticleId), int(in.ClassifyId)); err != nil {
		return nil, err
	}

	return &web.ArticleMoveResponse{}, nil
}

func (a *Article) Asterisk(ctx context.Context, in *web.ArticleAsteriskRequest) (*web.ArticleAsteriskResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)
	if err := a.ArticleService.Asterisk(ctx, uid, int(in.ArticleId), int(in.Action)); err != nil {
		return nil, err
	}

	return &web.ArticleAsteriskResponse{}, nil
}

func (a *Article) SetTags(ctx context.Context, in *web.ArticleTagsRequest) (*web.ArticleTagsResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := a.ArticleService.Tag(ctx, uid, int(in.ArticleId), in.GetTagIds()); err != nil {
		return nil, err
	}

	return &web.ArticleTagsResponse{}, nil
}

func (a *Article) RecoverList(ctx context.Context, _ *web.ArticleRecoverListRequest) (*web.ArticleRecoverListResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	items := make([]*web.ArticleRecoverListResponse_Item, 0)

	list, err := a.ArticleRepo.FindAll(ctx, func(db *gorm.DB) {
		db.Where("user_id = ? and status = ?", uid, 2)
		db.Where("deleted_at > ?", time.Now().Add(-time.Hour*24*30))
		db.Order("deleted_at desc,id desc")
	})

	if err != nil {
		return nil, err
	}

	classList, err := a.ArticleClassRepo.FindByIds(ctx, lo.Map(list, func(item *model.Article, index int) any {
		return item.ClassId
	}))

	if err != nil {
		return nil, err
	}

	classListMap := lo.KeyBy(classList, func(item *model.ArticleClass) int {
		return item.Id
	})

	for _, item := range list {
		className := ""

		if class, ok := classListMap[item.ClassId]; ok {
			className = class.ClassName
		}

		at := time.Until(item.DeletedAt.Time.Add(time.Hour * 24 * 30))

		items = append(items, &web.ArticleRecoverListResponse_Item{
			ArticleId:    int32(item.Id),
			ClassifyId:   int32(item.ClassId),
			ClassifyName: className,
			Title:        item.Title,
			Abstract:     item.Abstract,
			Image:        item.Image,
			CreatedAt:    item.CreatedAt.Format(time.DateTime),
			DeletedAt:    item.DeletedAt.Time.Format(time.DateTime),
			Day:          int32(math.Ceil(at.Seconds() / 86400)),
		})
	}

	return &web.ArticleRecoverListResponse{
		Items: items,
	}, nil
}
