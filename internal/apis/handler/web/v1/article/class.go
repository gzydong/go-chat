package article

import (
	"context"

	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/errorx"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/utils"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/service"
	"github.com/samber/lo"
)

var _ web.IArticleClassHandler = (*Class)(nil)

type Class struct {
	ArticleClassService service.IArticleClassService
}

func (c Class) List(ctx context.Context, req *web.ArticleClassListRequest) (*web.ArticleClassListResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)
	uid := session.GetAuthID()

	list, err := c.ArticleClassService.List(ctx, uid)
	if err != nil {
		return nil, err
	}

	items := make([]*web.ArticleClassListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ArticleClassListResponse_Item{
			Id:        int32(item.Id),
			ClassName: item.ClassName,
			IsDefault: int32(item.IsDefault),
			Count:     int32(item.Count),
		})
	}

	_, ok := lo.Find(list, func(item *model.ArticleClassItem) bool {
		return item.IsDefault == 1
	})

	if !ok {
		id, err := c.ArticleClassService.Create(ctx, uid, "默认分类", model.Yes)
		if err != nil {
			return nil, err
		}

		items = append(items, &web.ArticleClassListResponse_Item{
			Id:        int32(id),
			ClassName: "默认分类",
			IsDefault: model.Yes,
			Count:     0,
		})
	}

	return &web.ArticleClassListResponse{
		Items: items,
	}, nil
}

func (c Class) Edit(ctx context.Context, in *web.ArticleClassEditRequest) (*web.ArticleClassEditResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)
	uid := session.GetAuthID()

	if in.Name == "默认分类" {
		return nil, errorx.New(40001, "该分类名称禁止被创建/编辑")
	}

	if in.ClassifyId == 0 {
		id, err := c.ArticleClassService.Create(ctx, uid, in.Name, model.No)
		if err == nil {
			in.ClassifyId = int32(id)
		}
	} else {
		class, err := c.ArticleClassService.Find(ctx, int(in.ClassifyId))
		if err != nil {
			if utils.IsSqlNoRows(err) {
				return nil, entity.ErrNoteClassNotExist
			}

			return nil, err
		}

		if class.IsDefault == model.Yes {
			return nil, entity.ErrNoteClassDefaultNotAllow
		}

		err = c.ArticleClassService.Update(ctx, uid, int(in.ClassifyId), in.Name)
		if err != nil {
			return nil, err
		}
	}

	return &web.ArticleClassEditResponse{
		ClassifyId: in.ClassifyId,
	}, nil
}

func (c Class) Delete(ctx context.Context, in *web.ArticleClassDeleteRequest) (*web.ArticleClassDeleteResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)
	uid := session.GetAuthID()

	class, err := c.ArticleClassService.Find(ctx, int(in.ClassifyId))
	if err != nil {
		if utils.IsSqlNoRows(err) {
			return nil, entity.ErrNoteClassNotExist
		}

		return nil, err
	}

	if class.IsDefault == model.Yes {
		return nil, entity.ErrNoteClassDefaultNotDelete
	}

	err = c.ArticleClassService.Delete(ctx, uid, int(in.ClassifyId))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c Class) Sort(ctx context.Context, in *web.ArticleClassSortRequest) (*web.ArticleClassSortResponse, error) {
	session, _ := middleware.FormContext[entity.WebClaims](ctx)
	uid := session.UserId

	err := c.ArticleClassService.Sort(ctx, uid, in.ClassifyIds)
	if err != nil {
		return nil, err
	}

	return &web.ArticleClassSortResponse{}, nil
}
