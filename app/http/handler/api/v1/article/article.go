package article

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/service/note"
)

type Article struct {
	service *note.ArticleService
}

func NewArticleHandler(service *note.ArticleService) *Article {
	return &Article{service}
}

// List 文章列表
func (c *Article) List(ctx *gin.Context) {

}

// Detail 文章详情
func (c *Article) Detail(ctx *gin.Context) {

}

// Class 添加或编辑文章
func (c *Article) Edit(ctx *gin.Context) {

}

// Delete 删除文章
func (c *Article) Delete(ctx *gin.Context) {

}

// Recover 恢复文章
func (c *Article) Recover(ctx *gin.Context) {

}

// Upload 文章图片上传
func (c *Article) Upload(ctx *gin.Context) {

}

// Move 文章移动
func (c *Article) Move(ctx *gin.Context) {

}

// Asterisk 标记文章
func (c Article) Asterisk(ctx *gin.Context) {

}

// UpdateTag 文章标签
func (c *Article) UpdateTag(ctx *gin.Context) {

}

// ForeverDelete 永久删除文章
func (c *Article) ForeverDelete(ctx *gin.Context) {

}
