package article

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/service"
)

type Class struct {
	ArticleClassService service.IArticleClassService
}

// List 分类列表
func (c *Class) List(ctx *core.Context) error {

	list, err := c.ArticleClassService.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
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

	return ctx.Success(&web.ArticleClassListResponse{
		Items: items,
		Paginate: &web.Paginate{
			Page:  1,
			Size:  100000,
			Total: int32(len(items)),
		},
	})
}

// Edit 添加或修改分类
func (c *Class) Edit(ctx *core.Context) error {

	var (
		err error
		in  = &web.ArticleClassEditRequest{}
		uid = ctx.UserId()
	)

	if err = ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if in.ClassifyId == 0 {
		id, err := c.ArticleClassService.Create(ctx.Ctx(), uid, in.Name)
		if err == nil {
			in.ClassifyId = int32(id)
		}
	} else {
		err = c.ArticleClassService.Update(ctx.Ctx(), uid, int(in.ClassifyId), in.Name)
	}

	if err != nil {
		return ctx.ErrorBusiness("笔记分类编辑失败")
	}

	return ctx.Success(&web.ArticleClassEditResponse{
		ClassifyId: in.ClassifyId,
	})
}

// Delete 删除分类
func (c *Class) Delete(ctx *core.Context) error {

	in := &web.ArticleClassDeleteRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleClassService.Delete(ctx.Ctx(), ctx.UserId(), int(in.ClassifyId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleClassDeleteResponse{})
}

// Sort 删除分类
func (c *Class) Sort(ctx *core.Context) error {

	in := &web.ArticleClassSortRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleClassService.Sort(ctx.Ctx(), ctx.UserId(), int(in.ClassifyId), int(in.SortType))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleClassSortResponse{})
}
