package talk

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api2 "go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ginutil"
	"go-chat/internal/pkg/timeutil"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jwtutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/service"
)

type Records struct {
	service            *service.TalkRecordsService
	groupMemberService *service.GroupMemberService
	fileSystem         *filesystem.Filesystem
	authPermission     *service.AuthPermissionService
}

func NewTalkRecordsHandler(service *service.TalkRecordsService, groupMemberService *service.GroupMemberService, fileSystem *filesystem.Filesystem, authPermission *service.AuthPermissionService) *Records {
	return &Records{
		service:            service,
		groupMemberService: groupMemberService,
		fileSystem:         fileSystem,
		authPermission:     authPermission,
	}
}

// GetRecords 获取会话记录
func (c *Records) GetRecords(ctx *gin.Context) error {
	params := &api2.GetTalkRecordsRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)
	if params.TalkType == entity.ChatGroupMode {
		if !c.authPermission.IsAuth(ctx.Request.Context(), &service.AuthPermission{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}) {
			rows := make([]entity.H, 0)
			rows = append(rows, entity.H{
				"content":     "暂无权限查看群消息",
				"created_at":  timeutil.DateTime(),
				"id":          1,
				"msg_type":    0,
				"receiver_id": params.ReceiverId,
				"talk_type":   params.TalkType,
				"user_id":     0,
			})

			return ginutil.Success(ctx, entity.H{
				"limit":     params.Limit,
				"record_id": 0,
				"rows":      rows,
			})
		}
	}

	records, err := c.service.GetTalkRecords(ctx, &service.QueryTalkRecordsOpts{
		TalkType:   params.TalkType,
		UserId:     jwtutil.GetUid(ctx),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Id
	}

	return ginutil.Success(ctx, entity.H{
		"limit":     params.Limit,
		"record_id": rid,
		"rows":      records,
	})
}

// SearchHistoryRecords 查询下会话记录
func (c *Records) SearchHistoryRecords(ctx *gin.Context) error {
	params := &api2.GetTalkRecordsRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)

	if params.TalkType == entity.ChatGroupMode {
		if !c.authPermission.IsAuth(ctx.Request.Context(), &service.AuthPermission{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}) {
			return ginutil.Success(ctx, entity.H{
				"limit":     params.Limit,
				"record_id": 0,
				"rows":      make([]entity.H, 0),
			})
		}
	}

	m := []int{
		entity.MsgTypeText,
		entity.MsgTypeFile,
		entity.MsgTypeForward,
		entity.MsgTypeCode,
		entity.MsgTypeVote,
	}

	if sliceutil.InInt(params.MsgType, m) {
		m = []int{params.MsgType}
	}

	records, err := c.service.GetTalkRecords(ctx, &service.QueryTalkRecordsOpts{
		TalkType:   params.TalkType,
		MsgType:    m,
		UserId:     jwtutil.GetUid(ctx),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Id
	}

	return ginutil.Success(ctx, entity.H{
		"limit":     params.Limit,
		"record_id": rid,
		"rows":      records,
	})
}

// GetForwardRecords 获取转发记录
func (c *Records) GetForwardRecords(ctx *gin.Context) error {
	params := &api2.GetForwardTalkRecordRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	records, err := c.service.GetForwardRecords(ctx.Request.Context(), jwtutil.GetUid(ctx), int64(params.RecordId))
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	return ginutil.Success(ctx, entity.H{
		"rows": records,
	})
}

// Download 聊天文件下载
func (c *Records) Download(ctx *gin.Context) error {
	params := &api2.DownloadChatFileRequest{}
	if err := ctx.ShouldBindQuery(params); err != nil {
		return ginutil.InvalidParams(ctx, err)
	}

	resp, err := c.service.Dao().FindFileRecord(ctx.Request.Context(), params.RecordId)
	if err != nil {
		return ginutil.BusinessError(ctx, err)
	}

	uid := jwtutil.GetUid(ctx)
	if uid != resp.Record.UserId {
		if resp.Record.TalkType == entity.ChatPrivateMode {
			if resp.Record.ReceiverId != uid {
				return ginutil.Unauthorized(ctx, "无访问权限！")
			}
		} else {
			if !c.groupMemberService.Dao().IsMember(resp.Record.ReceiverId, uid, false) {
				return ginutil.Unauthorized(ctx, "无访问权限！")
			}
		}
	}

	switch resp.FileInfo.Drive {
	case entity.FileDriveLocal:
		ctx.FileAttachment(c.fileSystem.Local.Path(resp.FileInfo.Path), resp.FileInfo.OriginalName)
	case entity.FileDriveCos:
		ctx.Redirect(http.StatusFound, c.fileSystem.Cos.PrivateUrl(resp.FileInfo.Path, 60))
	default:
		return ginutil.BusinessError(ctx, "未知文件驱动类型")
	}

	return nil
}
