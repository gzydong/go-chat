package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/dao"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/filesystem"
	"net/http"
)

type Download struct {
	fileSystem     *filesystem.Filesystem
	talkRecordsDao *dao.TalkRecordsDao
}

func NewDownloadHandler(fileSystem *filesystem.Filesystem, talkRecordsDao *dao.TalkRecordsDao) *Download {
	return &Download{fileSystem, talkRecordsDao}
}

// ChatFile 下载聊天文件
func (d *Download) ChatFile(ctx *gin.Context) {
	params := &request.DownloadChatFileRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	resp, err := d.talkRecordsDao.FindFileRecord(ctx.Request.Context(), params.CrId)
	if err != nil {
		return
	}

	switch resp.FileInfo.SaveType {
	case 1:
		filePath := d.fileSystem.Local.Path(resp.FileInfo.SaveDir)
		ctx.FileAttachment(filePath, resp.FileInfo.OriginalName)
		return
	case 2:
		dwUrl := d.fileSystem.Cos.PrivateUrl(resp.FileInfo.SaveDir, 60)
		ctx.Redirect(http.StatusFound, dwUrl)
		return
	}
}

// ArticleAnnex 下载笔记附件
func (d *Download) ArticleAnnex(ctx *gin.Context) {

}
