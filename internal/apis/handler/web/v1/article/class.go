package article

import (
	"github.com/samber/lo"
	"go-chat/api/pb/web/v1"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type Class struct {
	ArticleClassService service.IArticleClassService
}

// List 分类列表
func (c *Class) List(ctx *core.Context) error {

	list, err := c.ArticleClassService.List(ctx.GetContext(), ctx.GetAuthId())
	if err != nil {
		return ctx.Error(err)
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
		id, err := c.ArticleClassService.Create(ctx.GetContext(), ctx.GetAuthId(), "默认分类", model.Yes)
		if err != nil {
			return ctx.Error(err)
		}

		items = append(items, &web.ArticleClassListResponse_Item{
			Id:        int32(id),
			ClassName: "默认分类",
			IsDefault: model.Yes,
			Count:     0,
		})
	}

	return ctx.Success(&web.ArticleClassListResponse{
		Items: items,
	})
}

// Edit 添加或修改分类
func (c *Class) Edit(ctx *core.Context) error {

	var (
		err error
		in  = &web.ArticleClassEditRequest{}
		uid = ctx.GetAuthId()
	)

	if err = ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if in.Name == "默认分类" {
		return ctx.InvalidParams("该分类名称禁止被创建/编辑")
	}

	if in.ClassifyId == 0 {
		id, err := c.ArticleClassService.Create(ctx.GetContext(), uid, in.Name, model.No)
		if err == nil {
			in.ClassifyId = int32(id)
		}
	} else {
		class, err := c.ArticleClassService.Find(ctx.GetContext(), int(in.ClassifyId))
		if err != nil {
			if utils.IsSqlNoRows(err) {
				return ctx.Error(entity.ErrNoteClassNotExist)
			}

			return ctx.Error(err)
		}

		if class.IsDefault == model.Yes {
			return ctx.Error(entity.ErrNoteClassDefaultNotAllow)
		}

		err = c.ArticleClassService.Update(ctx.GetContext(), uid, int(in.ClassifyId), in.Name)
		if err != nil {
			return ctx.Error(err)
		}
	}

	return ctx.Success(&web.ArticleClassEditResponse{
		ClassifyId: in.ClassifyId,
	})
}

// Delete 删除分类
func (c *Class) Delete(ctx *core.Context) error {

	in := &web.ArticleClassDeleteRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	class, err := c.ArticleClassService.Find(ctx.GetContext(), int(in.ClassifyId))
	if err != nil {
		if utils.IsSqlNoRows(err) {
			return ctx.Error(entity.ErrNoteClassNotExist)
		}

		return ctx.Error(err)
	}

	if class.IsDefault == model.Yes {
		return ctx.Error(entity.ErrNoteClassDefaultNotDelete)
	}

	err = c.ArticleClassService.Delete(ctx.GetContext(), ctx.GetAuthId(), int(in.ClassifyId))
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.ArticleClassDeleteResponse{})
}

// Sort 删除分类
func (c *Class) Sort(ctx *core.Context) error {

	in := &web.ArticleClassSortRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleClassService.Sort(ctx.GetContext(), ctx.GetAuthId(), int(in.ClassifyId), int(in.SortType))
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(&web.ArticleClassSortResponse{})
}
