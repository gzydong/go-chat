package article

import (
	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/service"
)

type Tag struct {
	ArticleTagService service.IArticleTagService
}

// List 标签列表
func (c *Tag) List(ctx *core.Context) error {

	list, err := c.ArticleTagService.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	items := make([]*web.ArticleTagListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ArticleTagListResponse_Item{
			Id:      int32(item.Id),
			TagName: item.TagName,
			Count:   int32(item.Count),
		})
	}

	return ctx.Success(&web.ArticleTagListResponse{Tags: items})
}

// Edit 添加或修改标签
func (c *Tag) Edit(ctx *core.Context) error {

	var (
		err error
		in  = &web.ArticleTagEditRequest{}
		uid = ctx.UserId()
	)

	if err = ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if in.TagId == 0 {
		id, err := c.ArticleTagService.Create(ctx.Ctx(), uid, in.TagName)
		if err == nil {
			in.TagId = int32(id)
		}
	} else {
		err = c.ArticleTagService.Update(ctx.Ctx(), uid, int(in.TagId), in.TagName)
	}

	if err != nil {
		return ctx.ErrorBusiness("笔记标签编辑失败")
	}

	return ctx.Success(&web.ArticleTagEditResponse{TagId: in.TagId})
}

// Delete 删除标签
func (c *Tag) Delete(ctx *core.Context) error {

	in := &web.ArticleTagDeleteRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleTagService.Delete(ctx.Ctx(), ctx.UserId(), int(in.TagId))
	if err != nil {
		return ctx.ErrorBusiness(err)
	}

	return ctx.Success(&web.ArticleTagDeleteResponse{})
}
