package talk

import (
	"net/http"
	"slices"
	"time"

	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/ichat"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Records struct {
	GroupMemberRepo    *repo.GroupMember
	TalkRecordsRepo    *repo.TalkRecords
	TalkRecordsService service.ITalkRecordsService
	GroupMemberService service.IGroupMemberService
	AuthService        service.IAuthService
	Filesystem         filesystem.IFilesystem
}

type GetTalkRecordsRequest struct {
	TalkType   int `form:"talk_type" json:"talk_type" binding:"required,oneof=1 2"`         // 对话类型
	MsgType    int `form:"msg_type" json:"msg_type" binding:"numeric"`                      // 消息类型
	ReceiverId int `form:"receiver_id" json:"receiver_id" binding:"required,numeric,min=1"` // 接收者ID
	Cursor     int `form:"cursor" json:"cursor" binding:"min=0,numeric"`                    // 上次查询的游标
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
		err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		})

		if err != nil {
			items := make([]entity.TalkRecord, 0)
			items = append(items, entity.TalkRecord{
				Id:         1,
				MsgId:      strutil.NewMsgId(),
				Sequence:   1,
				TalkType:   params.TalkType,
				MsgType:    entity.ChatMsgSysText,
				ReceiverId: params.ReceiverId,
				Extra: model.TalkRecordExtraText{
					Content: "暂无权限查看群消息",
				},
				CreatedAt: timeutil.DateTime(),
			})

			return ctx.Success(map[string]any{
				"cursor": 1,
				"items":  items,
			})
		}
	}

	records, err := c.TalkRecordsService.FindAllTalkRecords(ctx.Ctx(), &service.FindAllTalkRecordsOpt{
		TalkType:   params.TalkType,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
		Cursor:     params.Cursor,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	for i, record := range records {
		if record.IsRevoke == 1 {
			records[i].Extra = make(map[string]any)
		}
	}

	return ctx.Success(map[string]any{
		"cursor": cursor,
		"items":  records,
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
		err := c.AuthService.IsAuth(ctx.Ctx(), &service.AuthOption{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		})

		if err != nil {
			return ctx.Success(map[string]any{
				"cursor": 0,
				"items":  make([]map[string]any, 0),
			})
		}
	}

	m := []int{
		entity.ChatMsgTypeText,
		entity.ChatMsgTypeCode,
		entity.ChatMsgTypeImage,
		entity.ChatMsgTypeVideo,
		entity.ChatMsgTypeAudio,
		entity.ChatMsgTypeFile,
		entity.ChatMsgTypeLocation,
		entity.ChatMsgTypeForward,
		entity.ChatMsgTypeVote,
	}

	if slices.Contains(m, params.MsgType) {
		m = []int{params.MsgType}
	}

	records, err := c.TalkRecordsService.FindAllTalkRecords(ctx.Ctx(), &service.FindAllTalkRecordsOpt{
		TalkType:   params.TalkType,
		MsgType:    m,
		UserId:     ctx.UserId(),
		ReceiverId: params.ReceiverId,
		Cursor:     params.Cursor,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	for i, record := range records {
		if record.IsRevoke == 1 {
			records[i].Extra = make(map[string]any)
		}
	}

	return ctx.Success(map[string]any{
		"cursor": cursor,
		"items":  records,
	})
}

type GetForwardTalkRecordRequest struct {
	MsgId string `form:"msg_id" json:"msg_id" binding:"required"`
}

// GetForwardRecords 获取转发记录
func (c *Records) GetForwardRecords(ctx *ichat.Context) error {

	params := &GetForwardTalkRecordRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	records, err := c.TalkRecordsService.FindForwardRecords(ctx.Ctx(), ctx.UserId(), params.MsgId)
	if err != nil {
		return ctx.ErrorBusiness(err.Error())
	}

	return ctx.Success(map[string]any{
		"items": records,
	})
}

type DownloadChatFileRequest struct {
	MsgId string `form:"msg_id" json:"msg_id" binding:"required"`
}

// Download 聊天文件下载
func (c *Records) Download(ctx *ichat.Context) error {
	params := &DownloadChatFileRequest{}
	if err := ctx.Context.ShouldBindQuery(params); err != nil {
		return ctx.InvalidParams(err)
	}

	record, err := c.TalkRecordsRepo.FindByMsgId(ctx.Ctx(), params.MsgId)
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
			if !c.GroupMemberRepo.IsMember(ctx.Ctx(), record.ReceiverId, uid, false) {
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
		if c.Filesystem.Driver() != filesystem.LocalDriver {
			return ctx.ErrorBusiness("未知文件驱动类型")
		}

		filePath := c.Filesystem.(*filesystem.LocalFilesystem).Path(c.Filesystem.BucketPrivateName(), fileInfo.Path)
		ctx.Context.FileAttachment(filePath, fileInfo.Name)
	case entity.FileDriveMinio:
		ctx.Context.Redirect(http.StatusFound, c.Filesystem.PrivateUrl(c.Filesystem.BucketPrivateName(), fileInfo.Path, fileInfo.Name, 60*time.Second))
	default:
		return ctx.ErrorBusiness("未知文件驱动类型")
	}

	return nil
}
