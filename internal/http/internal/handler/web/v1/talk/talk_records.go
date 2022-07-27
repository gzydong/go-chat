package talk

import (
	"net/http"

	"go-chat/internal/http/internal/dto/web"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/timeutil"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/service"
)

type Records struct {
	service            *service.TalkRecordsService
	groupMemberService *service.GroupMemberService
	fileSystem         *filesystem.Filesystem
	authPermission     *service.AuthPermissionService
}

func NewRecords(service *service.TalkRecordsService, groupMemberService *service.GroupMemberService, fileSystem *filesystem.Filesystem, authPermission *service.AuthPermissionService) *Records {
	return &Records{service: service, groupMemberService: groupMemberService, fileSystem: fileSystem, authPermission: authPermission}
}

// GetRecords 获取会话记录
func (c *Records) GetRecords(ctx *ichat.Context) error {

	params := &web.GetTalkRecordsRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if params.TalkType == entity.ChatGroupMode {
		if !c.authPermission.IsAuth(ctx.RequestCtx(), &service.AuthPermission{
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

			return ctx.Success(entity.H{
				"limit":     params.Limit,
				"record_id": 0,
				"rows":      rows,
			})
		}
	}

	records, err := c.service.GetTalkRecords(ctx.Context, &service.QueryTalkRecordsOpt{
		TalkType:   params.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Id
	}

	return ctx.Success(entity.H{
		"limit":     params.Limit,
		"record_id": rid,
		"rows":      records,
	})
}

// SearchHistoryRecords 查询下会话记录
func (c *Records) SearchHistoryRecords(ctx *ichat.Context) error {

	params := &web.GetTalkRecordsRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if params.TalkType == entity.ChatGroupMode {
		if !c.authPermission.IsAuth(ctx.RequestCtx(), &service.AuthPermission{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}) {
			return ctx.Success(entity.H{
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

	records, err := c.service.GetTalkRecords(ctx.Context, &service.QueryTalkRecordsOpt{
		TalkType:   params.TalkType,
		MsgType:    m,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Id
	}

	return ctx.Success(entity.H{
		"limit":     params.Limit,
		"record_id": rid,
		"rows":      records,
	})
}

// GetForwardRecords 获取转发记录
func (c *Records) GetForwardRecords(ctx *ichat.Context) error {

	params := &web.GetForwardTalkRecordRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	records, err := c.service.GetForwardRecords(ctx.RequestCtx(), ctx.UserId(), int64(params.RecordId))
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	return ctx.Success(entity.H{
		"rows": records,
	})
}

// Download 聊天文件下载
func (c *Records) Download(ctx *ichat.Context) error {

	params := &web.DownloadChatFileRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	resp, err := c.service.Dao().FindFileRecord(ctx.RequestCtx(), params.RecordId)
	if err != nil {
		return ctx.BusinessError(err.Error())
	}

	uid := ctx.UserId()
	if uid != resp.Record.UserId {
		if resp.Record.TalkType == entity.ChatPrivateMode {
			if resp.Record.ReceiverId != uid {
				return ctx.Unauthorized("无访问权限！")
			}
		} else {
			if !c.groupMemberService.Dao().IsMember(resp.Record.ReceiverId, uid, false) {
				return ctx.Unauthorized("无访问权限！")
			}
		}
	}

	switch resp.FileInfo.Drive {
	case entity.FileDriveLocal:
		ctx.Context.FileAttachment(c.fileSystem.Local.Path(resp.FileInfo.Path), resp.FileInfo.OriginalName)
	case entity.FileDriveCos:
		ctx.Context.Redirect(http.StatusFound, c.fileSystem.Cos.PrivateUrl(resp.FileInfo.Path, 60))
	default:
		return ctx.BusinessError("未知文件驱动类型")
	}

	return nil
}
