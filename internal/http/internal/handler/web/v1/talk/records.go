package talk

import (
	"net/http"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/sliceutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/service"
)

type Records struct {
	talkRecordsService *service.TalkRecordsService
	groupMemberService *service.GroupMemberService
	filesystem         *filesystem.Filesystem
	authService        *service.AuthService
}

func NewRecords(talkRecordsService *service.TalkRecordsService, groupMemberService *service.GroupMemberService, filesystem *filesystem.Filesystem, authService *service.AuthService) *Records {
	return &Records{talkRecordsService: talkRecordsService, groupMemberService: groupMemberService, filesystem: filesystem, authService: authService}
}

type GetTalkRecordsRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2"`         // 对话类型
	MsgType    int `form:"msg_type" json:"msg_type" binding:"numeric"`                      // 消息类型
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,min=1"` // 接收者ID
	RecordId   int `form:"record_id" json:"record_id" binding:"min=0,numeric"`              // 上次查询的最小消息ID
	Limit      int `form:"limit" json:"limit" binding:"required,numeric,max=100"`           // 数据行数
}

// GetRecords 获取会话记录
func (c *Records) GetRecords(ctx *ichat.Context) error {

	params := &GetTalkRecordsRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()
	if params.TalkType == entity.ChatGroupMode {
		err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		})

		if err != nil {
			items := make([]map[string]any, 0)
			items = append(items, map[string]any{
				"content":     "暂无权限查看群消息",
				"created_at":  timeutil.DateTime(),
				"id":          1,
				"msg_id":      strutil.NewMsgId(),
				"msg_type":    0,
				"receiver_id": params.ReceiverId,
				"talk_type":   params.TalkType,
				"user_id":     0,
			})

			return ctx.Success(map[string]any{
				"limit":     params.Limit,
				"record_id": 0,
				"items":     items,
			})
		}
	}

	records, err := c.talkRecordsService.GetTalkRecords(ctx.Ctx(), &service.QueryTalkRecordsOpt{
		TalkType:   params.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Sequence
	}

	return ctx.Success(map[string]any{
		"limit":     params.Limit,
		"record_id": rid,
		"items":     records,
	})
}

// SearchHistoryRecords 查询下会话记录
func (c *Records) SearchHistoryRecords(ctx *ichat.Context) error {

	params := &GetTalkRecordsRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.UserId()

	if params.TalkType == entity.ChatGroupMode {
		err := c.authService.IsAuth(ctx.Ctx(), &service.AuthOption{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		})

		if err != nil {
			return ctx.Success(map[string]any{
				"limit":     params.Limit,
				"record_id": 0,
				"items":     make([]map[string]any, 0),
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

	if sliceutil.Include(params.MsgType, m) {
		m = []int{params.MsgType}
	}

	records, err := c.talkRecordsService.GetTalkRecords(ctx.Ctx(), &service.QueryTalkRecordsOpt{
		TalkType:   params.TalkType,
		MsgType:    m,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
		RecordId:   params.RecordId,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	rid := 0
	if length := len(records); length > 0 {
		rid = records[length-1].Sequence
	}

	return ctx.Success(map[string]any{
		"limit":     params.Limit,
		"record_id": rid,
		"items":     records,
	})
}

type GetForwardTalkRecordRequest struct {
	RecordId int `form:"record_id" json:"record_id" binding:"min=0,numeric"` // 上次查询的最小消息ID
}

// GetForwardRecords 获取转发记录
func (c *Records) GetForwardRecords(ctx *ichat.Context) error {

	params := &GetForwardTalkRecordRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	records, err := c.talkRecordsService.GetForwardRecords(ctx.Ctx(), ctx.UserId(), int64(params.RecordId))
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(map[string]any{
		"items": records,
	})
}

type DownloadChatFileRequest struct {
	RecordId int `form:"cr_id" json:"cr_id" binding:"required,min=1"`
}

// Download 聊天文件下载
func (c *Records) Download(ctx *ichat.Context) error {

	params := &DownloadChatFileRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	record, err := c.talkRecordsService.Dao().FindById(ctx.Ctx(), params.RecordId)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	uid := ctx.UserId()
	if uid != record.UserId {
		if record.TalkType == entity.ChatPrivateMode {
			if record.ReceiverId != uid {
				return ctx.Forbidden("无访问权限！")
			}
		} else {
			if !c.groupMemberService.Dao().IsMember(ctx.Ctx(), record.ReceiverId, uid, false) {
				return ctx.Forbidden("无访问权限！")
			}
		}
	}

	var fileInfo model.TalkRecordExtraFile
	if err := jsonutil.Decode(record.Extra, &fileInfo); err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	switch fileInfo.Drive {
	case entity.FileDriveLocal:
		ctx.Context.FileAttachment(c.filesystem.Local.Path(fileInfo.Path), fileInfo.OriginalName)
	case entity.FileDriveCos:
		ctx.Context.Redirect(http.StatusFound, c.filesystem.Cos.PrivateUrl(fileInfo.Path, 60*time.Second))
	default:
		return ctx.ErrorBusiness("未知文件驱动类型")
	}

	return nil
}
