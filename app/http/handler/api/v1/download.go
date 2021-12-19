package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/pkg/auth"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/service"
	"net/http"
)

type Download struct {
	fileSystem         *filesystem.Filesystem
	talkRecordsDao     *dao.TalkRecordsDao
	groupMemberService *service.GroupMemberService
}

func NewDownloadHandler(fileSystem *filesystem.Filesystem, talkRecordsDao *dao.TalkRecordsDao, groupMemberService *service.GroupMemberService) *Download {
	return &Download{fileSystem, talkRecordsDao, groupMemberService}
}

// TalkFile 下载聊天文件
func (c *Download) TalkFile(ctx *gin.Context) {
	params := &request.DownloadChatFileRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	resp, err := c.talkRecordsDao.FindFileRecord(ctx.Request.Context(), params.RecordId)
	if err != nil {
		return
	}

	uid := auth.GetAuthUserID(ctx)
	if uid != resp.Record.UserId {
		if resp.Record.TalkType == entity.PrivateChat {
			if resp.Record.ReceiverId != uid {
				response.Unauthorized(ctx, "无访问权限！")
				return
			}
		} else {
			if !c.groupMemberService.Dao().IsMember(resp.Record.ReceiverId, uid) {
				response.Unauthorized(ctx, "无访问权限！")
				return
			}
		}
	}

	switch resp.FileInfo.SaveType {
	case 1:
		filePath := c.fileSystem.Local.Path(resp.FileInfo.SaveDir)
		ctx.FileAttachment(filePath, resp.FileInfo.OriginalName)
		return
	case 2:
		ctx.Redirect(http.StatusFound, c.fileSystem.Cos.PrivateUrl(resp.FileInfo.SaveDir, 60))
		return
	}
}

// ArticleAnnex 下载笔记附件
func (c *Download) ArticleAnnex(ctx *gin.Context) {

}
