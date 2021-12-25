package article

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/service/note"
)

type Annex struct {
	service *note.ArticleAnnexService
}

func NewAnnexHandler(service *note.ArticleAnnexService) *Annex {
	return &Annex{service}
}

// List 附件列表
func (c *Annex) List(ctx *gin.Context) {

}

// Upload 上传附件
func (c *Annex) Upload(ctx *gin.Context) {

}

// Delete 删除附件
func (c *Annex) Delete(ctx *gin.Context) {

}

// Recover 恢复附件
func (c *Annex) Recover(ctx *gin.Context) {

}

// RecoverList 附件回收站列表
func (c *Annex) RecoverList(ctx *gin.Context) {

}

// ForeverDelete 永久删除附件
func (c *Annex) ForeverDelete(ctx *gin.Context) {

}
