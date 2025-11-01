package talk

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gzydong/go-chat/api/pb/web/v1"
	"github.com/gzydong/go-chat/internal/entity"
	"github.com/gzydong/go-chat/internal/pkg/core/errorx"
	"github.com/gzydong/go-chat/internal/pkg/core/middleware"
	"github.com/gzydong/go-chat/internal/pkg/filesystem"
	"github.com/gzydong/go-chat/internal/pkg/jsonutil"
	"github.com/gzydong/go-chat/internal/pkg/strutil"
	"github.com/gzydong/go-chat/internal/pkg/timeutil"
	"github.com/gzydong/go-chat/internal/repository/model"
	"github.com/gzydong/go-chat/internal/repository/repo"
	"github.com/gzydong/go-chat/internal/service"
	"github.com/samber/lo"
)

var _ web.IMessageHandler = (*Message)(nil)

type Message struct {
	TalkService          service.ITalkService
	AuthService          service.IAuthService
	Filesystem           filesystem.IFilesystem
	GroupMemberRepo      *repo.GroupMember
	TalkRecordFriendRepo *repo.TalkUserMessage
	TalkRecordGroupRepo  *repo.TalkGroupMessage
	TalkRecordsService   service.ITalkRecordService
	GroupMemberService   service.IGroupMemberService
}

func (m *Message) Records(ctx context.Context, in *web.MessageRecordsRequest) (*web.MessageRecordsResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if in.TalkMode == entity.ChatGroupMode {
		err := m.AuthService.IsAuth(ctx, &service.AuthOption{
			TalkType: int(in.TalkMode),
			UserId:   uid,
			ToFromId: int(in.ToFromId),
		})

		if err != nil {
			return &web.MessageRecordsResponse{
				Items: []*web.MessageRecord{
					{
						MsgId:     strutil.NewMsgId(),
						Sequence:  1,
						MsgType:   entity.ChatMsgSysText,
						FromId:    0,
						IsRevoked: model.No,
						SendTime:  timeutil.DateTime(),
						Extra: jsonutil.Encode(model.TalkRecordExtraText{
							Content: "暂无权限查看群消息",
						}),
						Quote: "{}",
					},
				},
				Cursor: 1,
			}, nil
		}
	}

	records, err := m.TalkRecordsService.FindAllTalkRecords(ctx, &service.FindAllTalkRecordsOpt{
		TalkType:   int(in.TalkMode),
		UserId:     uid,
		ReceiverId: int(in.ToFromId),
		Cursor:     int(in.Cursor),
		Limit:      int(in.Limit),
	})

	if err != nil {
		return nil, err
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	return &web.MessageRecordsResponse{
		Items: lo.Map(records, func(item *model.TalkMessageRecord, _ int) *web.MessageRecord {
			return &web.MessageRecord{
				FromId:    int32(item.FromId),
				MsgId:     item.MsgId,
				Sequence:  int32(item.Sequence),
				MsgType:   int32(item.MsgType),
				Nickname:  item.Nickname,
				Avatar:    item.Avatar,
				IsRevoked: int32(item.IsRevoked),
				SendTime:  item.SendTime.Format(time.DateTime),
				Extra:     lo.Ternary(item.IsRevoked == model.Yes, "{}", item.Extra),
				Quote:     item.Quote,
			}
		}),
		Cursor: int32(cursor),
	}, nil
}

func (m *Message) HistoryRecords(ctx context.Context, in *web.MessageHistoryRecordsRequest) (*web.MessageHistoryRecordsResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if in.TalkMode == entity.ChatGroupMode {
		err := m.AuthService.IsAuth(ctx, &service.AuthOption{
			TalkType: int(in.TalkMode),
			UserId:   uid,
			ToFromId: int(in.ToFromId),
		})

		if err != nil {
			return &web.MessageHistoryRecordsResponse{}, nil
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

	if slices.Contains(msgTypes, int(in.MsgType)) {
		msgTypes = []int{int(in.MsgType)}
	}

	records, err := m.TalkRecordsService.FindAllTalkRecords(ctx, &service.FindAllTalkRecordsOpt{
		TalkType:   int(in.TalkMode),
		MsgType:    msgTypes,
		UserId:     uid,
		ReceiverId: int(in.ToFromId),
		Cursor:     int(in.Cursor),
		Limit:      int(in.Limit),
	})

	if err != nil {
		return nil, err
	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	return &web.MessageHistoryRecordsResponse{
		Items: lo.Map(records, func(item *model.TalkMessageRecord, _ int) *web.MessageRecord {
			return &web.MessageRecord{
				FromId:    int32(item.FromId),
				MsgId:     item.MsgId,
				Sequence:  int32(item.Sequence),
				MsgType:   int32(item.MsgType),
				Nickname:  item.Nickname,
				Avatar:    item.Avatar,
				IsRevoked: int32(item.IsRevoked),
				SendTime:  item.SendTime.Format(time.DateTime),
				Extra:     lo.Ternary(item.IsRevoked == model.Yes, "{}", item.Extra),
				Quote:     item.Quote,
			}
		}),
		Cursor: int32(cursor),
	}, nil
}

func (m *Message) ForwardRecords(ctx context.Context, in *web.MessageForwardRecordsRequest) (*web.MessageRecordsClearResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	records, err := m.TalkRecordsService.FindForwardRecords(ctx, uid, in.MsgIds, int(in.TalkMode))
	if err != nil {
		return nil, err
	}

	return &web.MessageRecordsClearResponse{
		Items: lo.Map(records, func(item *model.TalkMessageRecord, _ int) *web.MessageRecord {
			return &web.MessageRecord{
				FromId:    int32(item.FromId),
				MsgId:     item.MsgId,
				Sequence:  int32(item.Sequence),
				MsgType:   int32(item.MsgType),
				Nickname:  item.Nickname,
				Avatar:    item.Avatar,
				IsRevoked: int32(item.IsRevoked),
				SendTime:  item.SendTime.Format(time.DateTime),
				Extra:     lo.Ternary(item.IsRevoked == model.Yes, "{}", item.Extra),
				Quote:     item.Quote,
			}
		}),
	}, nil
}

func (m *Message) Revoke(ctx context.Context, in *web.MessageRevokeRequest) (*web.MessageRevokeResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := m.TalkService.Revoke(ctx, &service.TalkRevokeOption{
		UserId:   uid,
		TalkMode: int(in.TalkMode),
		MsgId:    in.MsgId,
	}); err != nil {
		return nil, err
	}

	return &web.MessageRevokeResponse{}, nil
}

func (m *Message) Delete(ctx context.Context, in *web.MessageDeleteRequest) (*web.MessageDeleteResponse, error) {
	uid := middleware.FormContextAuthId[entity.WebClaims](ctx)

	if err := m.TalkService.DeleteRecord(ctx, &service.TalkDeleteRecordOption{
		UserId:   uid,
		TalkMode: int(in.TalkMode),
		ToFromId: int(in.ToFromId),
		MsgIds:   in.MsgIds,
	}); err != nil {
		return nil, err
	}

	return &web.MessageDeleteResponse{}, nil
}

type DownloadChatFileRequest struct {
	TalkMode int    `form:"talk_mode" json:"talk_mode" binding:"required,oneof=1 2"`
	MsgId    string `form:"msg_id" json:"msg_id" binding:"required"`
}

// Download 聊天文件下载
func (m *Message) Download(ctx *gin.Context) error {
	params := &DownloadChatFileRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		return errorx.New(400, err.Error())
	}

	uid := middleware.FormContextAuthId[entity.WebClaims](ctx.Request.Context())

	var fileInfo model.TalkRecordExtraFile
	if params.TalkMode == entity.ChatGroupMode {
		record, err := m.TalkRecordGroupRepo.FindByWhere(ctx, "msg_id = ?", params.MsgId)
		if err != nil {
			return ctx.Error(err)
		}

		if !m.GroupMemberRepo.IsMember(ctx, record.GroupId, uid, false) {
			return entity.ErrPermissionDenied
		}

		if err := jsonutil.Unmarshal(record.Extra, &fileInfo); err != nil {
			return err
		}
	} else {
		record, err := m.TalkRecordFriendRepo.FindByWhere(ctx, "user_id = ? and msg_id = ?", uid, params.MsgId)
		if err != nil {
			return errorx.New(400, "文件不存在")
		}

		if err := jsonutil.Unmarshal(record.Extra, &fileInfo); err != nil {
			return err
		}
	}

	switch m.Filesystem.Driver() {
	case filesystem.LocalDriver:
		filePath := m.Filesystem.(*filesystem.LocalFilesystem).Path(m.Filesystem.BucketPrivateName(), fileInfo.Path)
		ctx.FileAttachment(filePath, fileInfo.Name)
	case filesystem.MinioDriver:
		ctx.Redirect(http.StatusFound, m.Filesystem.PrivateUrl(m.Filesystem.BucketPrivateName(), fileInfo.Path, fileInfo.Name, 60*time.Second))
	default:
		return errorx.New(400, "未知文件驱动类型")
	}

	return nil
}
