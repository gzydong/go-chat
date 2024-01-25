package message

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/business"
	"time"

	"github.com/google/uuid"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

var _ IService = (*Service)(nil)

// IPrivateMessage 私有消息
type IPrivateMessage interface {
	// CreatePrivateSysMessage 给指定用户创建私有的系统消息
	CreatePrivateSysMessage(ctx context.Context, option CreatePrivateSysMessageOption) error
	// CreatePrivateMessage 创建私有的消息
	CreatePrivateMessage(ctx context.Context, option CreatePrivateMessageOption) error
	// CreateToUserPrivateMessage 给指定用户信箱添加消息
	CreateToUserPrivateMessage(ctx context.Context, data *model.TalkUserMessage) error
}

// IGroupMessage 群消息
type IGroupMessage interface {
	// CreateGroupMessage 创建群消息
	CreateGroupMessage(ctx context.Context, option CreateGroupMessageOption) error
	// CreateGroupSysMessage 创建群系统消息
	CreateGroupSysMessage(ctx context.Context, option CreateGroupSysMessageOption) error
}

type IMessage interface {
	// CreateMessage 创建消息
	CreateMessage(ctx context.Context, option CreateMessageOption) error
	// CreateLoginMessage 创建登录消息
	CreateLoginMessage(ctx context.Context, option CreateLoginMessageOption) error
	// CreateTextMessage 文本消息
	CreateTextMessage(ctx context.Context, option CreateTextMessage) error
	// CreateImageMessage 图片文件消息
	CreateImageMessage(ctx context.Context, option CreateImageMessage) error
	// CreateVoiceMessage 语音文件消息
	CreateVoiceMessage(ctx context.Context, option CreateVoiceMessage) error
	// CreateVideoMessage 视频文件消息
	CreateVideoMessage(ctx context.Context, option CreateVideoMessage) error
	// CreateFileMessage 文件消息
	CreateFileMessage(ctx context.Context, option CreateFileMessage) error
	// CreateCodeMessage 代码消息
	CreateCodeMessage(ctx context.Context, option CreateCodeMessage) error
	// CreateVoteMessage 投票消息
	CreateVoteMessage(ctx context.Context, option CreateVoteMessage) error
	// CreateEmoticonMessage 表情消息
	CreateEmoticonMessage(ctx context.Context, option CreateEmoticonMessage) error
	// CreateForwardMessage 转发消息
	CreateForwardMessage(ctx context.Context, option CreateForwardMessage) error
	// CreateLocationMessage 位置消息
	CreateLocationMessage(ctx context.Context, option CreateLocationMessage) error
	// CreateBusinessCardMessage 推送用户名片消息
	CreateBusinessCardMessage(ctx context.Context, option CreateBusinessCardMessage) error
	// CreateMixedMessage 图文消息
	CreateMixedMessage(ctx context.Context, option CreateMixedMessage) error
}

type IService interface {
	IPrivateMessage
	IGroupMessage
	IMessage
}

type Service struct {
	*repo.Source
	GroupMemberRepo     *repo.GroupMember
	SplitUploadRepo     *repo.FileUpload
	TalkRecordsVoteRepo *repo.GroupVote
	UsersRepo           *repo.Users
	Filesystem          filesystem.IFilesystem
	UnreadStorage       *cache.UnreadStorage
	MessageStorage      *cache.MessageStorage
	ServerStorage       *cache.ServerStorage
	ClientStorage       *cache.ClientStorage
	Sequence            *repo.Sequence
	RobotRepo           *repo.Robot

	PushMessage *business.PushMessage
}

func (s *Service) CreateMessage(ctx context.Context, option CreateMessageOption) error {
	if option.TalkMode == 1 {
		return s.CreatePrivateMessage(ctx, CreatePrivateMessageOption{
			MsgType:  option.MsgType,
			FromId:   option.FromId,
			ToFromId: option.ToFromId,
			QuoteId:  option.QuoteId,
			Extra:    option.Extra,
		})
	}

	return s.CreateGroupMessage(ctx, CreateGroupMessageOption{
		MsgType:  option.MsgType,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		QuoteId:  option.QuoteId,
		Extra:    option.Extra,
	})
}

func (s *Service) CreateTextMessage(ctx context.Context, option CreateTextMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeText,
		QuoteId:  option.QuoteId,
		Extra: jsonutil.Encode(model.TalkRecordExtraText{
			Content:  option.Content,
			Mentions: option.Mentions,
		}),
	})
}

func (s *Service) CreateImageMessage(ctx context.Context, option CreateImageMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeImage,
		QuoteId:  option.QuoteId,
		Extra: jsonutil.Encode(model.TalkRecordExtraImage{
			Size:   option.Size,
			Url:    option.Url,
			Width:  option.Width,
			Height: option.Height,
		}),
	})
}

func (s *Service) CreateVoiceMessage(ctx context.Context, option CreateVoiceMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeAudio,
		Extra: jsonutil.Encode(model.TalkRecordExtraAudio{
			Name:     "",
			Size:     option.Size,
			Url:      option.Url,
			Duration: option.Duration,
		}),
	})
}

func (s *Service) CreateVideoMessage(ctx context.Context, option CreateVideoMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeVideo,
		Extra: jsonutil.Encode(model.TalkRecordExtraVideo{
			Name:     "",
			Cover:    option.Cover,
			Size:     option.Size,
			Url:      option.Url,
			Duration: option.Duration,
		}),
	})
}

func (s *Service) CreateFileMessage(ctx context.Context, option CreateFileMessage) error {
	now := time.Now()

	file, err := s.SplitUploadRepo.GetFile(ctx, option.FromId, option.UploadId)
	if err != nil {
		return err
	}

	publicUrl := ""
	filePath := fmt.Sprintf("talk-files/%s/%s.%s", now.Format("200601"), uuid.New().String(), file.FileExt)

	// 公开文件
	if entity.GetMediaType(file.FileExt) <= 3 {
		filePath = strutil.GenMediaObjectName(file.FileExt, 0, 0)
		// 如果是多媒体文件，则将私有文件转移到公开文件
		if err := s.Filesystem.CopyObject(
			s.Filesystem.BucketPrivateName(), file.Path,
			s.Filesystem.BucketPublicName(), filePath,
		); err != nil {
			return err
		}

		publicUrl = s.Filesystem.PublicUrl(s.Filesystem.BucketPublicName(), filePath)
	} else {
		if err := s.Filesystem.Copy(s.Filesystem.BucketPrivateName(), file.Path, filePath); err != nil {
			return err
		}
	}

	message := CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
	}

	switch entity.GetMediaType(file.FileExt) {
	case entity.MediaFileAudio:
		message.MsgType = entity.ChatMsgTypeAudio
		message.Extra = jsonutil.Encode(&model.TalkRecordExtraAudio{
			Size:     int(file.FileSize),
			Url:      publicUrl,
			Duration: 0,
		})
	case entity.MediaFileVideo:
		message.MsgType = entity.ChatMsgTypeVideo
		message.Extra = jsonutil.Encode(&model.TalkRecordExtraVideo{
			Cover:    "",
			Size:     int(file.FileSize),
			Url:      publicUrl,
			Duration: 0,
		})
	case entity.MediaFileOther:
		message.MsgType = entity.ChatMsgTypeFile
		message.Extra = jsonutil.Encode(&model.TalkRecordExtraFile{
			Name: file.OriginalName,
			Size: int(file.FileSize),
			Path: filePath,
		})
	}

	return s.CreateMessage(ctx, message)
}

func (s *Service) CreateCodeMessage(ctx context.Context, option CreateCodeMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeCode,
		Extra: jsonutil.Encode(model.TalkRecordExtraCode{
			Lang: option.Lang,
			Code: option.Code,
		}),
	})
}

func (s *Service) CreateVoteMessage(ctx context.Context, option CreateVoteMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeVote,
		Extra: jsonutil.Encode(model.TalkRecordExtraVote{
			VoteId: option.VoteId,
		}),
	})
}

func (s *Service) CreateEmoticonMessage(ctx context.Context, option CreateEmoticonMessage) error {
	var emoticon model.EmoticonItem
	if err := s.Source.Db().First(&emoticon, "id = ? and user_id = ?", option.EmoticonId, option.FromId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("表情信息不存在")
		}

		return err
	}

	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeImage,
		Extra: jsonutil.Encode(model.TalkRecordExtraImage{
			Url: emoticon.Url,
		}),
	})
}

// CreateForwardMessage todo 待完善
func (s *Service) CreateForwardMessage(ctx context.Context, option CreateForwardMessage) error {
	items := make([]ForwardMessageOpt, 0)

	// 发送方式 1:逐条发送 2:合并发送
	if option.Action == 1 {
		for _, userId := range option.Uids {
			item := ForwardMessageOpt{
				MsgIds:       option.MsgIds,
				TalkMode:     option.TalkMode,
				ToFromId:     option.ToFromId,
				UserId:       option.FromId,
				ToUserId:     userId,
				ToUserIdType: 1,
			}

			items = append(items, item)

			err := s.toSplitForward(ctx, item)
			if err != nil {
				logger.WithFields(
					logger.LevelError,
					fmt.Sprintf("[uid] split forward message failed err:%s", err.Error()),
					item,
				)
			}
		}

		for _, groupId := range option.Gids {
			item := ForwardMessageOpt{
				MsgIds:       option.MsgIds,
				TalkMode:     option.TalkMode,
				ToFromId:     option.ToFromId,
				UserId:       option.FromId,
				ToUserId:     groupId,
				ToUserIdType: 2,
			}

			items = append(items, item)

			err := s.toSplitForward(ctx, item)
			if err != nil {
				logger.WithFields(
					logger.LevelError,
					fmt.Sprintf("[group] split forward message failed err:%s", err.Error()),
					item,
				)
			}
		}
	} else {
		for _, userId := range option.Uids {
			item := ForwardMessageOpt{
				MsgIds:   option.MsgIds,
				TalkMode: option.TalkMode,
				ToFromId: option.ToFromId,

				UserId:       option.UserId,
				ToUserId:     userId,
				ToUserIdType: 1,
			}

			items = append(items, item)

			err := s.toCombineForward(ctx, item)
			if err != nil {
				logger.WithFields(
					logger.LevelError,
					fmt.Sprintf("[uid] combin forward message failed err:%s", err.Error()),
					item,
				)
			}
		}

		for _, groupId := range option.Gids {
			item := ForwardMessageOpt{
				MsgIds:       option.MsgIds,
				TalkMode:     option.TalkMode,
				ToFromId:     option.ToFromId,
				UserId:       option.UserId,
				ToUserId:     groupId,
				ToUserIdType: 2,
			}

			items = append(items, item)

			err := s.toCombineForward(ctx, item)
			if err != nil {
				logger.WithFields(
					logger.LevelError,
					fmt.Sprintf("[group] combin forward message failed err:%s", err.Error()),
					item,
				)
			}
		}
	}

	return nil
}

func (s *Service) CreateLocationMessage(ctx context.Context, option CreateLocationMessage) error {
	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeLocation,
		Extra: jsonutil.Encode(model.TalkRecordExtraLocation{
			Longitude:   option.Longitude,
			Latitude:    option.Latitude,
			Description: option.Description,
		}),
	})
}

func (s *Service) CreateBusinessCardMessage(ctx context.Context, option CreateBusinessCardMessage) error {
	userInfo, err := s.UsersRepo.FindById(ctx, option.ToFromId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if userInfo == nil {
		return errors.New("用户不存在")
	}

	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeCard,
		Extra: jsonutil.Encode(model.TalkRecordExtraUserShare{
			UserId:   option.UserId,
			Nickname: userInfo.Nickname,
			Avatar:   userInfo.Avatar,
			Describe: userInfo.Motto,
		}),
	})
}

func (s *Service) CreateMixedMessage(ctx context.Context, option CreateMixedMessage) error {
	items := make([]*model.TalkRecordExtraMixedItem, 0)
	for _, item := range option.MessageList {
		items = append(items, &model.TalkRecordExtraMixedItem{
			Type:    item.Type,
			Content: item.Content,
		})
	}

	return s.CreateMessage(ctx, CreateMessageOption{
		TalkMode: option.TalkMode,
		FromId:   option.FromId,
		ToFromId: option.ToFromId,
		MsgType:  entity.ChatMsgTypeMixed,
		Extra: jsonutil.Encode(model.TalkRecordExtraMixed{
			Items: items,
		}),
	})
}

func (s *Service) CreateLoginMessage(ctx context.Context, option CreateLoginMessageOption) error {
	robot, err := s.RobotRepo.GetLoginRobot(ctx)
	if err != nil {
		return err
	}

	return s.CreateToUserPrivateMessage(ctx, &model.TalkUserMessage{
		MsgType:  entity.ChatMsgTypeLogin,
		UserId:   option.UserId,
		ToFromId: robot.UserId,
		FromId:   robot.UserId,
		Extra: jsonutil.Encode(&model.TalkRecordExtraLogin{
			IP:       option.Ip,
			Platform: option.Platform,
			Agent:    option.Agent,
			Address:  option.Address,
			Reason:   option.Reason,
			Datetime: option.LoginAt,
		}),
	})
}

func (s *Service) getTextMessage(msgType int, extra string) string {
	return text(msgType, extra)
}
