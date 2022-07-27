package article

import (
	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/service/note"
)

type Class struct {
	service *note.ArticleClassService
}

func NewClass(service *note.ArticleClassService) *Class {
	return &Class{service}
}

// List 分类列表
func (c *Class) List(ctx *ichat.Context) error {

	items, err := c.service.List(ctx.RequestCtx(), ctx.UserId())
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Paginate(items, 1, 100000, len(items))
}

// Edit 添加或修改分类
func (c *Class) Edit(ctx *ichat.Context) error {

	var (
		err    error
		params = &web.ArticleClassEditRequest{}
		uid    = ctx.UserId()
	)

	if err = ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.ClassId == 0 {
		params.ClassId, err = c.service.Create(ctx.RequestCtx(), uid, params.ClassName)
	} else {
		err = c.service.Update(ctx.RequestCtx(), uid, params.ClassId, params.ClassName)
	}

	if err != nil {
		return ctx.BusinessError("笔记分类编辑失败")
	}

	return ctx.Success(entity.H{"id": params.ClassId})
}

// Delete 删除分类
func (c *Class) Delete(ctx *ichat.Context) error {

	params := &web.ArticleClassDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.Delete(ctx.RequestCtx(), ctx.UserId(), params.ClassId)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Sort 删除分类
func (c *Class) Sort(ctx *ichat.Context) error {

	params := &web.ArticleClassSortRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.Sort(ctx.RequestCtx(), ctx.UserId(), params.ClassId, params.SortType)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}
