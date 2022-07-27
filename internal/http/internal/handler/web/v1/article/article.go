package article

import (
	"bytes"
	"fmt"

	"go-chat/internal/entity"
	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/repository/model"

	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/pkg/utils"
	"go-chat/internal/service/note"
)

type Article struct {
	service             *note.ArticleService
	fileSystem          *filesystem.Filesystem
	articleAnnexService *note.ArticleAnnexService
}

func NewArticle(service *note.ArticleService, fileSystem *filesystem.Filesystem, articleAnnexService *note.ArticleAnnexService) *Article {
	return &Article{service, fileSystem, articleAnnexService}
}

// List 文章列表
func (c *Article) List(ctx *ichat.Context) error {

	params := &web.ArticleListRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	items, err := c.service.List(ctx.RequestCtx(), &note.ArticleListOpt{
		UserId:   ctx.UserId(),
		Keyword:  params.Keyword,
		FindType: params.FindType,
		Cid:      params.Cid,
		Page:     params.Page,
	})
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	list := make([]map[string]interface{}, 0)
	for _, item := range items {
		list = append(list, map[string]interface{}{
			"id":          item.Id,
			"class_id":    item.ClassId,
			"tags_id":     item.TagsId,
			"title":       item.Title,
			"abstract":    item.Abstract,
			"class_name":  item.ClassName,
			"image":       item.Image,
			"is_asterisk": item.IsAsterisk,
			"status":      item.Status,
			"created_at":  timeutil.FormatDatetime(item.CreatedAt),
			"updated_at":  timeutil.FormatDatetime(item.UpdatedAt),
		})
	}

	return ctx.Paginate(list, 1, 1000, len(items))
}

// Detail 文章详情
func (c *Article) Detail(ctx *ichat.Context) error {

	params := &web.ArticleDetailRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	detail, err := c.service.Detail(ctx.RequestCtx(), uid, params.ArticleId)
	if err != nil {
		return ctx.BusinessError("笔记不存在")
	}

	tags := make([]map[string]interface{}, 0)
	for _, tagId := range sliceutil.ParseIds(detail.TagsId) {
		tags = append(tags, map[string]interface{}{"id": tagId})
	}

	files := make([]map[string]interface{}, 0)
	items, err := c.articleAnnexService.Dao().AnnexList(ctx.Context, uid, params.ArticleId)
	if err == nil {
		for _, item := range items {
			files = append(files, map[string]interface{}{
				"id":            item.Id,
				"suffix":        item.Suffix,
				"size":          item.Size,
				"original_name": item.OriginalName,
				"created_at":    timeutil.FormatDatetime(item.CreatedAt),
			})
		}
	}

	return ctx.Success(entity.H{
		"id":          detail.Id,
		"class_id":    detail.ClassId,
		"title":       detail.Title,
		"md_content":  detail.MdContent,
		"content":     detail.Content,
		"is_asterisk": detail.IsAsterisk,
		"status":      detail.Status,
		"created_at":  timeutil.FormatDatetime(detail.CreatedAt),
		"updated_at":  timeutil.FormatDatetime(detail.UpdatedAt),
		"tags":        tags,
		"files":       files,
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

	opts := &note.ArticleEditOpt{
		UserId:    uid,
		ArticleId: params.ArticleId,
		ClassId:   params.ClassId,
		Title:     params.Title,
		Content:   params.Content,
		MdContent: params.MdContent,
	}

	if params.ArticleId == 0 {
		params.ArticleId, err = c.service.Create(ctx.RequestCtx(), opts)
	} else {
		err = c.service.Update(ctx.RequestCtx(), opts)
	}

	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	var info *model.Article
	_ = c.service.Db().First(&info, params.ArticleId)

	return ctx.Success(entity.H{
		"id":       info.Id,
		"image":    info.Image,
		"abstract": info.Abstract,
		"title":    info.Title,
	})
}

// Delete 删除文章
func (c *Article) Delete(ctx *ichat.Context) error {

	params := &web.ArticleDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.UpdateStatus(ctx.RequestCtx(), ctx.UserId(), params.ArticleId, 2)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Recover 恢复文章
func (c *Article) Recover(ctx *ichat.Context) error {

	params := &web.ArticleRecoverRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	err := c.service.UpdateStatus(ctx.RequestCtx(), ctx.UserId(), params.ArticleId, 1)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Upload 文章图片上传
func (c *Article) Upload(ctx *ichat.Context) error {

	file, err := ctx.Context.FormFile("image")
	if err != nil {
		return ctx.InvalidParams("image 字段必传！")
	}

	if !sliceutil.InStr(strutil.FileSuffix(file.Filename), []string{"png", "jpg", "jpeg", "gif", "webp"}) {
		return ctx.InvalidParams("上传文件格式不正确,仅支持 png、jpg、jpeg、gif 和 webp")
	}

	// 判断上传文件大小（5M）
	if file.Size > 5<<20 {
		return ctx.InvalidParams("上传文件大小不能超过5M！")
	}

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return ctx.BusinessError("文件上传失败")
	}

	ext := strutil.FileSuffix(file.Filename)
	meta := utils.LoadImage(bytes.NewReader(stream))

	filePath := fmt.Sprintf("public/media/image/note/%s/%s", timeutil.DateNumber(), strutil.GenImageName(ext, meta.Width, meta.Height))

	if err := c.fileSystem.Default.Write(stream, filePath); err != nil {
		return ctx.BusinessError("文件上传失败")
	}

	return ctx.Success(entity.H{"url": c.fileSystem.Default.PublicUrl(filePath)})
}

// Move 文章移动
func (c *Article) Move(ctx *ichat.Context) error {

	params := &web.ArticleMoveRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.Move(ctx.RequestCtx(), ctx.UserId(), params.ArticleId, params.ClassId); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Asterisk 标记文章
func (c Article) Asterisk(ctx *ichat.Context) error {

	params := &web.ArticleAsteriskRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.Asterisk(ctx.RequestCtx(), ctx.UserId(), params.ArticleId, params.Type); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// Tag 文章标签
func (c *Article) Tag(ctx *ichat.Context) error {

	params := &web.ArticleTagsRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.Tag(ctx.RequestCtx(), ctx.UserId(), params.ArticleId, params.Tags); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}

// ForeverDelete 永久删除文章
func (c *Article) ForeverDelete(ctx *ichat.Context) error {

	params := &web.ArticleForeverDeleteRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	if err := c.service.ForeverDelete(ctx.RequestCtx(), ctx.UserId(), params.ArticleId); err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(nil)
}
