package article

import (
	"go-chat/api/pb/web/v1"
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

	list, err := c.service.List(ctx.Ctx(), ctx.UserId())
	if err != nil {
		return ctx.BusinessError(err)
	}

	items := make([]*web.ArticleTagListResponse_Item, 0, len(list))
	for _, item := range list {
		items = append(items, &web.ArticleTagListResponse_Item{
			Id:      int32(item.Id),
			TagName: item.TagName,
			Count:   int32(item.Count),
		})
	}

	return ctx.Success(&web.ArticleTagListResponse{
		Tags: items,
	})
}

// Edit 添加或修改标签
func (c *Tag) Edit(ctx *ichat.Context) error {

	var (
		err    error
		params = &web.ArticleTagEditRequest{}
		uid    = ctx.UserId()
	)

	if err = ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if params.TagId == 0 {
		id, err := c.service.Create(ctx.Ctx(), uid, params.TagName)
		if err == nil {
			params.TagId = int32(id)
		}
	} else {
		err = c.service.Update(ctx.Ctx(), uid, int(params.TagId), params.TagName)
	}

	if err != nil {
		return ctx.BusinessError("笔记标签编辑失败")
	}

	return ctx.Success(web.ArticleTagEditResponse{Id: params.TagId})
}

// Delete 删除标签
func (c *Tag) Delete(ctx *ichat.Context) error {

	params := &web.ArticleTagDeleteRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.Delete(ctx.Ctx(), ctx.UserId(), int(params.TagId))
	if err != nil {
		return ctx.BusinessError(err)
	}

	return ctx.Success(&web.ArticleTagDeleteResponse{})
}
