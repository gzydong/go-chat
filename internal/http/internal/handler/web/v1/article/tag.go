package article

import (
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service/note"
)

type Tag struct {
	service *note.ArticleTagService
}

func NewTag(service *note.ArticleTagService) *Tag {
	return &Tag{service}
}

// List 标签列表
func (c *Tag) List(ctx *ichat.Context) error {

	items, err := c.service.List(ctx.RequestCtx(), ctx.UserId())
	if err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(entity.H{"tags": items})
}

// Edit 添加或修改标签
func (c *Tag) Edit(ctx *ichat.Context) error {

	var (
		err    error
		params = &web.ArticleTagEditRequest{}
		uid    = ctx.UserId()
	)

	if err = ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.TagId == 0 {
		params.TagId, err = c.service.Create(ctx.RequestCtx(), uid, params.TagName)
	} else {
		err = c.service.Update(ctx.RequestCtx(), uid, params.TagId, params.TagName)
	}

	if err != nil {
		return ctx.BusinessError("笔记标签编辑失败")
	}

	return ctx.Success(entity.H{"id": params.TagId})
}

// Delete 删除标签
func (c *Tag) Delete(ctx *ichat.Context) error {

	params := &web.ArticleTagDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.Delete(ctx.RequestCtx(), ctx.UserId(), params.TagId)
	if err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(nil)
}
