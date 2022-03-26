package article

import (
	"github.com/gin-gonic/gin"

	"go-chat/internal/http/internal/request"
	"go-chat/internal/http/internal/response"
	"go-chat/internal/pkg/jwt"
	"go-chat/internal/service/note"
)

type Tag struct {
	service *note.ArticleTagService
}

func NewTagHandler(service *note.ArticleTagService) *Tag {
	return &Tag{service}
}

// 标签列表
func (c *Tag) List(ctx *gin.Context) {
	items, err := c.service.List(ctx.Request.Context(), jwt.GetUid(ctx))
	if err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, gin.H{"tags": items})
	}
}

// Edit 添加或修改标签
func (c *Tag) Edit(ctx *gin.Context) {
	var (
		err    error
		params = &request.ArticleTagEditRequest{}
		uid    = jwt.GetUid(ctx)
	)

	if err = ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if params.TagId == 0 {
		params.TagId, err = c.service.Create(ctx.Request.Context(), uid, params.TagName)
	} else {
		err = c.service.Update(ctx.Request.Context(), uid, params.TagId, params.TagName)
	}

	if err != nil {
		response.BusinessError(ctx, "笔记标签编辑失败")
	} else {
		response.Success(ctx, gin.H{"id": params.TagId})
	}
}

// 删除标签
func (c *Tag) Delete(ctx *gin.Context) {
	params := &request.ArticleTagDeleteRequest{}

	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	err := c.service.Delete(ctx.Request.Context(), jwt.GetUid(ctx), params.TagId)
	if err != nil {
		response.BusinessError(ctx, err)
	} else {
		response.Success(ctx, nil, "删除成功")
	}
}