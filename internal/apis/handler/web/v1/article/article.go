package article

import (
	"html"
	"time"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
	"gorm.io/gorm"
)

type Article struct {
	Source              *repo.Source
	ArticleAnnexRepo    *repo.ArticleAnnex
	ArticleRepo         *repo.Article
	ArticleService      service.IArticleService
	ArticleAnnexService service.IArticleAnnexService
	Filesystem          filesystem.IFilesystem
}

// List 文章列表
func (c *Article) List(ctx *core.Context) error {
	in := &web.ArticleListRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	items, err := c.ArticleService.List(ctx.Ctx(), &service.ArticleListOpt{
		UserId:     ctx.UserId(),
		FindType:   int(in.FindType),
		Keyword:    in.Keyword,
		ClassifyId: int(in.ClassifyId),
		TagId:      int(in.TagId),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	list := make([]*web.ArticleListResponse_Item, 0)
	for _, item := range items {
		list = append(list, &web.ArticleListResponse_Item{
			ArticleId:  int32(item.Id),
			ClassifyId: int32(item.ClassId),
			TagsId:     item.TagsId,
			Title:      item.Title,
			ClassName:  item.ClassName,
			Image:      item.Image,
			IsAsterisk: int32(item.IsAsterisk),
			Status:     int32(item.Status),
			CreatedAt:  timeutil.FormatDatetime(item.CreatedAt),
			UpdatedAt:  timeutil.FormatDatetime(item.UpdatedAt),
			Abstract:   item.Abstract,
		})
	}

	return ctx.Success(&web.ArticleListResponse{
		Items: list,
		Paginate: &web.ArticleListResponse_Paginate{
			Page:  1,
			Size:  1000,
			Total: int32(len(list)),
		},
	})
}

// Detail 文章详情
func (c *Article) Detail(ctx *core.Context) error {

	in := &web.ArticleDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	detail, err := c.ArticleService.Detail(ctx.Ctx(), uid, int(in.ArticleId))
	if err != nil {
		return ctx.ErrorBusiness("笔记不存在")
	}

	tags := make([]*web.ArticleDetailResponse_Tag, 0)
	for _, id := range sliceutil.ParseIds(detail.TagsId) {
		tags = append(tags, &web.ArticleDetailResponse_Tag{Id: int32(id)})
	}

	files := make([]*web.ArticleDetailResponse_AnnexFile, 0)
	items, err := c.ArticleAnnexRepo.AnnexList(ctx.Ctx(), uid, int(in.ArticleId))
	if err == nil {
		for _, item := range items {
			files = append(files, &web.ArticleDetailResponse_AnnexFile{
				AnnexId:   int32(item.Id),
				AnnexName: item.OriginalName,
				AnnexSize: int32(item.Size),
				CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
			})
		}
	}

	return ctx.Success(&web.ArticleDetailResponse{
		ArticleId:  int32(detail.Id),
		ClassifyId: int32(detail.ClassId),
		Title:      detail.Title,
		MdContent:  html.UnescapeString(detail.MdContent),
		IsAsterisk: int32(detail.IsAsterisk),
		CreatedAt:  timeutil.FormatDatetime(detail.CreatedAt),
		UpdatedAt:  timeutil.FormatDatetime(detail.UpdatedAt),
		TagIds:     tags,
		AnnexList:  files,
	})
}

// Editor 添加或编辑文章
func (c *Article) Editor(ctx *core.Context) error {

	var (
		err error
		in  = &web.ArticleEditRequest{}
		uid = ctx.UserId()
	)

	if err = ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	opt := &service.ArticleEditOpt{
		UserId:    uid,
		ArticleId: int(in.ArticleId),
		ClassId:   int(in.ClassifyId),
		Title:     in.Title,
		MdContent: in.MdContent,
	}

	if in.ArticleId == 0 {
		id, err := c.ArticleService.Create(ctx.Ctx(), opt)
		if err == nil {
			in.ArticleId = int32(id)
		}
	} else {
		err = c.ArticleService.Update(ctx.Ctx(), opt)
	}

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	var info *model.Article
	if err := c.Source.Db().First(&info, in.ArticleId).Error; err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleEditResponse{
		ArticleId: int32(info.Id),
		Title:     info.Title,
		Abstract:  info.Abstract,
		Image:     info.Image,
	})
}

// Delete 删除文章
func (c *Article) Delete(ctx *core.Context) error {

	in := &web.ArticleDeleteRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(in.ArticleId), 2)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(web.ArticleDeleteResponse{})
}

// Recover 恢复文章
func (c *Article) Recover(ctx *core.Context) error {

	in := &web.ArticleRecoverRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(in.ArticleId), 1)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleRecoverResponse{})
}

// MoveClassify 文章移动
func (c *Article) MoveClassify(ctx *core.Context) error {
	in := &web.ArticleMoveRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.Move(ctx.Ctx(), ctx.UserId(), int(in.ArticleId), int(in.ClassifyId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleMoveResponse{})
}

// Collect 标记文章
func (c *Article) Collect(ctx *core.Context) error {

	in := &web.ArticleAsteriskRequest{}
	if err := ctx.Context.ShouldBindJSON(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.Asterisk(ctx.Ctx(), ctx.UserId(), int(in.ArticleId), int(in.Action)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAsteriskResponse{})
}

// UpdateTag 文章标签
func (c *Article) UpdateTag(ctx *core.Context) error {

	in := &web.ArticleTagsRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.Tag(ctx.Ctx(), ctx.UserId(), int(in.ArticleId), in.GetTagIds()); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleTagsResponse{})
}

// ForeverDelete 永久删除文章
func (c *Article) ForeverDelete(ctx *core.Context) error {

	in := &web.ArticleForeverDeleteRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.ForeverDelete(ctx.Ctx(), ctx.UserId(), int(in.ArticleId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleForeverDeleteResponse{})
}

// RecycleList 永久删除文章
func (c *Article) RecycleList(ctx *core.Context) error {
	in := &web.ArticleRecoverListRequest{}
	if err := ctx.Context.ShouldBind(in); err != nil {
		return ctx.InvalidParams(err)
	}

	items := make([]*web.ArticleRecoverListResponse_Item, 0)

	list, err := c.ArticleRepo.FindAll(ctx.Ctx(), func(db *gorm.DB) {
		db.Where("user_id = ? and status = ?", ctx.UserId(), 2)
		db.Where("updated_at > ?", time.Now().Add(-time.Hour*24*7))
		db.Order("updated_at desc,id desc")
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	for _, item := range list {
		items = append(items, &web.ArticleRecoverListResponse_Item{
			ArticleId:    int32(item.Id),
			ClassifyId:   int32(item.ClassId),
			ClassifyName: "mlmlkmlkm",
			Title:        item.Title,
			Abstract:     item.Abstract,
			Image:        item.Image,
			CreatedAt:    item.CreatedAt.Format(time.DateTime),
			DeletedAt:    item.DeletedAt.Time.Format(time.DateTime),
		})
	}

	return ctx.Success(&web.ArticleRecoverListResponse{
		Items: items,
	})
}
