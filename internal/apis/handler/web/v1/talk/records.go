package talk

import (
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/samber/lo"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/core"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"go-chat/internal/service"
)

type Records struct {
	GroupMemberRepo      *repo.GroupMember
	TalkRecordFriendRepo *repo.TalkUserMessage
	TalkRecordGroupRepo  *repo.TalkGroupMessage
	TalkRecordsService   service.ITalkRecordService
	GroupMemberService   service.IGroupMemberService
	AuthService          service.IAuthService
	Filesystem           filesystem.IFilesystem
}

type GetTalkRecordsRequest struct {
	TalkMode int `form:"talk_mode" json:"talk_mode" binding:"required,oneof=1 2"`       // 对话类型
	ToFromId int `form:"to_from_id" json:"to_from_id" binding:"required,numeric,min=1"` // 接收者ID
	MsgType  int `form:"msg_type" json:"msg_type" binding:"numeric"`                    // 消息类型
	Cursor   int `form:"cursor" json:"cursor" binding:"min=0,numeric"`                  // 上次查询的游标
	Limit    int `form:"limit" json:"limit" binding:"required,numeric,max=100"`         // 数据行数
}

// GetRecords 获取会话记录
func (c *Records) GetRecords(ctx *core.Context) error {
	in := &GetTalkRecordsRequest{}
	if err := ctx.ShouldBindProto(in); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	if in.TalkMode == entity.ChatGroupMode {
		err := c.AuthService.IsAuth(ctx.GetContext(), &service.AuthOption{
			TalkType: in.TalkMode,
			UserId:   uid,
			ToFromId: in.ToFromId,
		})

		if err != nil {
			return ctx.Success(map[string]any{
				"cursor": 1,
				"items": []entity.ImMessagePayloadBody{
					{
						MsgId:     strutil.NewMsgId(),
						Sequence:  1,
						MsgType:   entity.ChatMsgSysText,
						FromId:    0,
						Nickname:  "",
						Avatar:    "",
						IsRevoked: model.No,
						SendTime:  timeutil.DateTime(),
						Extra: jsonutil.Encode(model.TalkRecordExtraText{
							Content: "暂无权限查看群消息",
						}),
						Quote: nil,
					},
				},
			})
		}
	}

	records, err := c.TalkRecordsService.FindAllTalkRecords(ctx.GetContext(), &service.FindAllTalkRecordsOpt{
		TalkType:   in.TalkMode,
		UserId:     uid,
		ReceiverId: in.ToFromId,
		Cursor:     in.Cursor,
		Limit:      in.Limit,
	})

	if err != nil {
		return ctx.Error(err)
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	return ctx.Success(map[string]any{
		"cursor": cursor,
		"items": lo.Map(records, func(item *model.TalkMessageRecord, index int) entity.ImMessagePayloadBody {
			return entity.ImMessagePayloadBody{
				FromId:    item.FromId,
				MsgId:     item.MsgId,
				Sequence:  item.Sequence,
				MsgType:   item.MsgType,
				Nickname:  item.Nickname,
				Avatar:    item.Avatar,
				IsRevoked: item.IsRevoked,
				SendTime:  item.SendTime.Format(time.DateTime),
				Extra:     lo.Ternary(item.IsRevoked == model.Yes, "{}", item.Extra),
				Quote:     item.Quote,
			}
		}),
	})
}

// SearchHistoryRecords 查询下会话记录
func (c *Records) SearchHistoryRecords(ctx *core.Context) error {

	params := &GetTalkRecordsRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	uid := ctx.AuthId()

	if params.TalkMode == entity.ChatGroupMode {
		err := c.AuthService.IsAuth(ctx.GetContext(), &service.AuthOption{
			TalkType: params.TalkMode,
			UserId:   uid,
			ToFromId: params.ToFromId,
		})

		if err != nil {
			return ctx.Success(map[string]any{
				"cursor": 0,
				"items":  make([]map[string]any, 0),
			})
		}
	}

	msgTypes := []int{
		entity.ChatMsgTypeText,
		entity.ChatMsgTypeMixed,
		entity.ChatMsgTypeCode,
		entity.ChatMsgTypeImage,
		entity.ChatMsgTypeVideo,
		entity.ChatMsgTypeAudio,
		entity.ChatMsgTypeFile,
		entity.ChatMsgTypeLocation,
		entity.ChatMsgTypeForward,
		entity.ChatMsgTypeVote,
	}

	if slices.Contains(msgTypes, params.MsgType) {
		msgTypes = []int{params.MsgType}
	}

	records, err := c.TalkRecordsService.FindAllTalkRecords(ctx.GetContext(), &service.FindAllTalkRecordsOpt{
		TalkType:   params.TalkMode,
		MsgType:    msgTypes,
		UserId:     uid,
		ReceiverId: params.ToFromId,
		Cursor:     params.Cursor,
		Limit:      params.Limit,
	})

	if err != nil {
		return ctx.Error(err)
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	return ctx.Success(map[string]any{
		"cursor": cursor,
		"items": lo.Map(records, func(item *model.TalkMessageRecord, index int) entity.ImMessagePayloadBody {
			return entity.ImMessagePayloadBody{
				FromId:    item.FromId,
				MsgId:     item.MsgId,
				Sequence:  item.Sequence,
				MsgType:   item.MsgType,
				Nickname:  item.Nickname,
				Avatar:    item.Avatar,
				IsRevoked: item.IsRevoked,
				SendTime:  item.SendTime.Format(time.DateTime),
				Extra:     lo.Ternary(item.IsRevoked == model.Yes, "{}", item.Extra),
				Quote:     item.Quote,
			}
		}),
	})
}

type GetForwardTalkRecordRequest struct {
	TalkMode int      `form:"talk_mode" json:"talk_mode" binding:"required,oneof=1 2"` // 对话类型
	MsgIds   []string `form:"msg_ids[]" json:"msg_ids" binding:"required"`
}

// GetForwardRecords 获取转发记录
func (c *Records) GetForwardRecords(ctx *core.Context) error {
	params := &GetForwardTalkRecordRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	records, err := c.TalkRecordsService.FindForwardRecords(ctx.GetContext(), ctx.AuthId(), params.MsgIds, params.TalkMode)
	if err != nil {
		return ctx.Error(err)
	}

	return ctx.Success(map[string]any{
		"items": lo.Map(records, func(item *model.TalkMessageRecord, index int) entity.ImMessagePayloadBody {
			return entity.ImMessagePayloadBody{
				FromId:    item.FromId,
				MsgId:     item.MsgId,
				Sequence:  item.Sequence,
				MsgType:   item.MsgType,
				Nickname:  item.Nickname,
				Avatar:    item.Avatar,
				IsRevoked: item.IsRevoked,
				SendTime:  item.SendTime.Format(time.DateTime),
				Extra:     lo.Ternary(item.IsRevoked == model.Yes, "{}", item.Extra),
				Quote:     item.Quote,
			}
		}),
	})
}

type DownloadChatFileRequest struct {
	TalkMode int    `form:"talk_mode" json:"talk_mode" binding:"required,oneof=1 2"`
	MsgId    string `form:"msg_id" json:"msg_id" binding:"required"`
}

// Download 聊天文件下载
func (c *Records) Download(ctx *core.Context) error {
	params := &DownloadChatFileRequest{}
	if err := ctx.Context.ShouldBind(params); err != nil {
		return ctx.InvalidParams(err)
	}

	var fileInfo model.TalkRecordExtraFile
	if params.TalkMode == entity.ChatGroupMode {
		record, err := c.TalkRecordGroupRepo.FindByWhere(ctx.GetContext(), "msg_id = ?", params.MsgId)
		if err != nil {
			return ctx.Error(err)
		}

		if !c.GroupMemberRepo.IsMember(ctx.GetContext(), record.GroupId, ctx.AuthId(), false) {
			return ctx.Forbidden("无访问权限！")
		}

		if err := jsonutil.Unmarshal(record.Extra, &fileInfo); err != nil {
			return ctx.Error(err)
		}
	} else {
		record, err := c.TalkRecordFriendRepo.FindByWhere(ctx.GetContext(), "user_id = ? and msg_id = ?", ctx.AuthId(), params.MsgId)
		if err != nil {
			return ctx.Error(err)
		}

		if err := jsonutil.Unmarshal(record.Extra, &fileInfo); err != nil {
			return ctx.Error(err)
		}
	}

	switch c.Filesystem.Driver() {
	case filesystem.LocalDriver:
		filePath := c.Filesystem.(*filesystem.LocalFilesystem).Path(c.Filesystem.BucketPrivateName(), fileInfo.Path)
		ctx.Context.FileAttachment(filePath, fileInfo.Name)
	case filesystem.MinioDriver:
		ctx.Context.Redirect(http.StatusFound, c.Filesystem.PrivateUrl(c.Filesystem.BucketPrivateName(), fileInfo.Path, fileInfo.Name, 60*time.Second))
	default:
		return ctx.Error(errors.New("未知文件驱动类型"))
	}

	return nil
}
