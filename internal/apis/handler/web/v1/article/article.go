package article

import (
	"bytes"
	"html"
	"slices"

	"go-chat/api/pb/web/v1"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Article struct {
	Source              *repo.Source
	ArticleAnnexRepo    *repo.ArticleAnnex
	ArticleService      service.IArticleService
	ArticleAnnexService service.IArticleAnnexService
	Filesystem          filesystem.IFilesystem
}

// List 文章列表
func (c *Article) List(ctx *ichat.Context) error {

	params := &web.ArticleListRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	items, err := c.ArticleService.List(ctx.Ctx(), &service.ArticleListOpt{
		UserId:   ctx.UserId(),
		Keyword:  params.Keyword,
		FindType: int(params.FindType),
		Cid:      int(params.Cid),
		Page:     int(params.Page),
	})
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	list := make([]*web.ArticleListResponse_Item, 0)
	for _, item := range items {
		list = append(list, &web.ArticleListResponse_Item{
			Id:         int32(item.Id),
			ClassId:    int32(item.ClassId),
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
func (c *Article) Detail(ctx *ichat.Context) error {

	params := &web.ArticleDetailRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	detail, err := c.ArticleService.Detail(ctx.Ctx(), uid, int(params.ArticleId))
	if err != nil {
		return ctx.ErrorBusiness("笔记不存在")
	}

	tags := make([]*web.ArticleDetailResponse_Tag, 0)
	for _, id := range sliceutil.ParseIds(detail.TagsId) {
		tags = append(tags, &web.ArticleDetailResponse_Tag{Id: int32(id)})
	}

	files := make([]*web.ArticleDetailResponse_File, 0)
	items, err := c.ArticleAnnexRepo.AnnexList(ctx.Ctx(), uid, int(params.ArticleId))
	if err == nil {
		for _, item := range items {
			files = append(files, &web.ArticleDetailResponse_File{
				Id:           int32(item.Id),
				Suffix:       item.Suffix,
				Size:         int32(item.Size),
				OriginalName: item.OriginalName,
				CreatedAt:    timeutil.FormatDatetime(item.CreatedAt),
			})
		}
	}

	return ctx.Success(&web.ArticleDetailResponse{
		Id:         int32(detail.Id),
		ClassId:    int32(detail.ClassId),
		Title:      detail.Title,
		MdContent:  html.UnescapeString(detail.MdContent),
		IsAsterisk: int32(detail.IsAsterisk),
		CreatedAt:  timeutil.FormatDatetime(detail.CreatedAt),
		UpdatedAt:  timeutil.FormatDatetime(detail.UpdatedAt),
		Tags:       tags,
		Files:      files,
	})
}

// Edit 添加或编辑文章
func (c *Article) Edit(ctx *ichat.Context) error {

	var (
		err    error
		params = &web.ArticleEditRequest{}
		uid    = ctx.UserId()
	)

	if err = ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	opt := &service.ArticleEditOpt{
		UserId:    uid,
		ArticleId: int(params.ArticleId),
		ClassId:   int(params.ClassId),
		Title:     params.Title,
		MdContent: params.MdContent,
	}

	if params.ArticleId == 0 {
		id, err := c.ArticleService.Create(ctx.Ctx(), opt)
		if err == nil {
			params.ArticleId = int32(id)
		}
	} else {
		err = c.ArticleService.Update(ctx.Ctx(), opt)
	}

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	var info *model.Article
	if err := c.Source.Db().First(&info, params.ArticleId).Error; err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleEditResponse{
		Id:       int32(info.Id),
		Title:    info.Title,
		Abstract: info.Abstract,
		Image:    info.Image,
	})
}

// Delete 删除文章
func (c *Article) Delete(ctx *ichat.Context) error {

	params := &web.ArticleDeleteRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(params.ArticleId), 2)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(web.ArticleDeleteResponse{})
}

// Recover 恢复文章
func (c *Article) Recover(ctx *ichat.Context) error {

	params := &web.ArticleRecoverRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.ArticleService.UpdateStatus(ctx.Ctx(), ctx.UserId(), int(params.ArticleId), 1)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleRecoverResponse{})
}

// Upload 文章图片上传
func (c *Article) Upload(ctx *ichat.Context) error {

	file, err := ctx.Context.FormFile("image")
	if err != nil {
		return ctx.InvalidParams("image 字段必传！")
	}

	if !slices.Contains([]string{"png", "jpg", "jpeg", "gif", "webp"}, strutil.FileSuffix(file.Filename)) {
		return ctx.InvalidParams("上传文件格式不正确,仅支持 png、jpg、jpeg、gif 和 webp")
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		return ctx.InvalidParams("上传文件大小不能超过5M！")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return ctx.ErrorBusiness("文件上传失败")
	}

	ext := strutil.FileSuffix(file.Filename)
	meta := utils.ReadImageMeta(bytes.NewReader(stream))

	filePath := strutil.GenMediaObjectName(ext, meta.Width, meta.Height)
	if err := c.Filesystem.Write(c.Filesystem.BucketPublicName(), filePath, stream); err != nil {
		return ctx.ErrorBusiness("文件上传失败")
	}

	return ctx.Success(map[string]any{"url": c.Filesystem.PublicUrl(c.Filesystem.BucketPublicName(), filePath)})
}

// Move 文章移动
func (c *Article) Move(ctx *ichat.Context) error {

	params := &web.ArticleMoveRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.Move(ctx.Ctx(), ctx.UserId(), int(params.ArticleId), int(params.ClassId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleMoveResponse{})
}

// Asterisk 标记文章
func (c *Article) Asterisk(ctx *ichat.Context) error {

	params := &web.ArticleAsteriskRequest{}
	if err := ctx.Context.ShouldBindJSON(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.Asterisk(ctx.Ctx(), ctx.UserId(), int(params.ArticleId), int(params.Type)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleAsteriskResponse{})
}

// Tag 文章标签
func (c *Article) Tag(ctx *ichat.Context) error {

	params := &web.ArticleTagsRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.Tag(ctx.Ctx(), ctx.UserId(), int(params.ArticleId), params.Tags); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleTagsResponse{})
}

// ForeverDelete 永久删除文章
func (c *Article) ForeverDelete(ctx *ichat.Context) error {

	params := &web.ArticleForeverDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.ArticleService.ForeverDelete(ctx.Ctx(), ctx.UserId(), int(params.ArticleId)); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(&web.ArticleForeverDeleteResponse{})
}
