package article

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/service/note"
)

type Class struct {
	service *note.ArticleClassService
}

func NewClassHandler(service *note.ArticleClassService) *Class {
	return &Class{service}
}

// List 分类列表
func (c *Class) List(ctx *gin.Context) {
	items, err := c.service.List(ctx.Request.Context(), auth.GetAuthUserID(ctx))
	if err != nil {
		response.BusinessError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"rows": items,
	})
}

// Edit 添加或修改分类
func (c *Class) Edit(ctx *gin.Context) {
	var (
		err    error
		params = &request.ArticleClassEditRequest{}
		uid    = auth.GetAuthUserID(ctx)
	)

	if err = ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if params.ClassId == 0 {
		params.ClassId, err = c.service.Create(ctx.Request.Context(), uid, params.ClassName)
	} else {
		err = c.service.Update(ctx.Request.Context(), uid, params.ClassId, params.ClassName)
	}

	if err != nil {
		response.BusinessError(ctx, "笔记标签编辑失败")
		return
	}

	response.Success(ctx, gin.H{
		"id": params.ClassId,
	})
}

// Delete 删除分类
func (c *Class) Delete(ctx *gin.Context) {
	params := &request.ArticleClassDeleteRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Delete(ctx.Request.Context(), auth.GetAuthUserID(ctx), params.ClassId)
	if err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}

// Sort 删除分类
func (c *Class) Sort(ctx *gin.Context) {
	params := &request.ArticleClassSortRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Sort(ctx.Request.Context(), auth.GetAuthUserID(ctx), params.ClassId, params.SortType)
	if err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil, "操作成功")
	}
}
