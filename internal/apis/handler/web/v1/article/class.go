package article

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service"
)

type Class struct {
	ArticleClassService service.IArticleClassService
}

// List 分类列表
func (c *Class) List(ctx *ichat.Context) error {

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
func (c *Class) Edit(ctx *ichat.Context) error {

	var (
		err    error
		params = &web.ArticleClassEditRequest{}
		uid    = ctx.UserId()
	)

	if err = ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.ClassId == 0 {
		id, err := c.ArticleClassService.Create(ctx.Ctx(), uid, params.ClassName)
		if err == nil {
			params.ClassId = int32(id)
		}
	} else {
		err = c.ArticleClassService.Update(ctx.Ctx(), uid, int(params.ClassId), params.ClassName)
	}

	if err != nil {
		return ctx.ErrorBusiness("笔记分类编辑失败")
	}

	return ctx.Success(&web.ArticleClassEditResponse{
		Id: params.ClassId,
	})
}

// Delete 删除分类
func (c *Class) Delete(ctx *ichat.Context) error {

	params := &web.ArticleClassDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleClassService.Delete(ctx.Ctx(), ctx.UserId(), int(params.ClassId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleClassDeleteResponse{})
}

// Sort 删除分类
func (c *Class) Sort(ctx *ichat.Context) error {

	params := &web.ArticleClassSortRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleClassService.Sort(ctx.Ctx(), ctx.UserId(), int(params.ClassId), int(params.SortType))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleClassSortResponse{})
}
