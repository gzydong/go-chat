package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/url"
	"strconv"

	"github.com/sirupsen/logrus"
	"go-chat/api/pb/message/v1"
	"go-chat/internal/entity"
	"go-chat/internal/logic"
	"go-chat/internal/pkg/encrypt"
	"go-chat/internal/pkg/filesystem"
	"go-chat/internal/pkg/jsonutil"
	"go-chat/internal/pkg/logger"
	"go-chat/internal/pkg/strutil"
	"go-chat/internal/pkg/timeutil"
	"go-chat/internal/repository/cache"
	"go-chat/internal/repository/model"
	"go-chat/internal/repository/repo"
	"gorm.io/gorm"
)

type MessageService struct {
	*repo.Source
	forward         *logic.MessageForwardLogic
	groupMemberRepo *repo.GroupMember
	splitUploadRepo *repo.SplitUpload
	fileSystem      *filesystem.Filesystem
	unreadStorage   *cache.UnreadStorage
	messageStorage  *cache.MessageStorage
	sidStorage      *cache.ServerStorage
	clientStorage   *cache.ClientStorage
	Sequence        *repo.Sequence
}

func NewMessageService(source *repo.Source, forward *logic.MessageForwardLogic, groupMemberRepo *repo.GroupMember, splitUploadRepo *repo.SplitUpload, fileSystem *filesystem.Filesystem, unreadStorage *cache.UnreadStorage, messageStorage *cache.MessageStorage, sidStorage *cache.ServerStorage, clientStorage *cache.ClientStorage, sequence *repo.Sequence) *MessageService {
	return &MessageService{Source: source, forward: forward, groupMemberRepo: groupMemberRepo, splitUploadRepo: splitUploadRepo, fileSystem: fileSystem, unreadStorage: unreadStorage, messageStorage: messageStorage, sidStorage: sidStorage, clientStorage: clientStorage, Sequence: sequence}
}

// SendText 文本消息
func (m *MessageService) SendText(ctx context.Context, uid int, req *message.TextMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeText,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Content:    html.EscapeString(req.Content),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	if err := m.Db().WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	m.afterHandle(ctx, data, map[string]string{
		"text": strutil.MtSubstr(data.Content, 0, 300),
	})

	return nil
}

// SendImage 图片文件消息
func (m *MessageService) SendImage(ctx context.Context, uid int, req *message.ImageMessageRequest) error {

	parse, err := url.Parse(req.Url)
	if err != nil {
		return err
	}

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeFile,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraFile{
			Type:         entity.MediaFileImage,
			Drive:        entity.FileDriveMode("local"),
			OriginalName: "图片名称",
			Suffix:       strutil.FileSuffix(req.Url),
			Size:         int(req.Size),
			Path:         parse.Path,
			Url:          req.Url,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err = m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[图片消息]"})
	}

	return err
}

// SendVoice 语音文件消息
func (m *MessageService) SendVoice(ctx context.Context, uid int, req *message.VoiceMessageRequest) error {

	parse, err := url.Parse(req.Url)
	if err != nil {
		return err
	}

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeFile,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraFile{
			Type:         entity.MediaFileAudio,
			Drive:        entity.FileDriveMode("local"),
			OriginalName: "语音文件",
			Suffix:       strutil.FileSuffix(req.Url),
			Size:         int(req.Size),
			Path:         parse.Path,
			Url:          req.Url,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err = m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[语音消息]"})
	}

	return err
}

// SendVideo 视频文件消息
func (m *MessageService) SendVideo(ctx context.Context, uid int, req *message.VideoMessageRequest) error {

	parse, err := url.Parse(req.Url)
	if err != nil {
		return err
	}

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeFile,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraFile{
			Type:         entity.MediaFileVideo,
			Drive:        entity.FileDriveMode("local"),
			OriginalName: "语音文件",
			Suffix:       strutil.FileSuffix(req.Url),
			Size:         int(req.Size),
			Path:         parse.Path,
			Url:          req.Url,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err = m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[文件消息]"})
	}

	return err
}

// SendFile 文件消息
func (m *MessageService) SendFile(ctx context.Context, uid int, req *message.FileMessageRequest) error {

	file, err := m.splitUploadRepo.GetFile(ctx, uid, req.UploadId)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("private/files/talks/%s/%s.%s", timeutil.DateNumber(), encrypt.Md5(strutil.Random(16)), file.FileExt)
	uri := ""
	if entity.GetMediaType(file.FileExt) <= 3 {
		filePath = fmt.Sprintf("public/media/%s/%s.%s", timeutil.DateNumber(), encrypt.Md5(strutil.Random(16)), file.FileExt)
		uri = m.fileSystem.Default.PublicUrl(filePath)
	}

	if err := m.fileSystem.Default.Copy(file.Path, filePath); err != nil {
		logrus.Error("文件拷贝失败 err: ", err.Error())
		return err
	}

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeFile,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraFile{
			Type:         entity.GetMediaType(file.FileExt),
			Drive:        file.Drive,
			OriginalName: file.OriginalName,
			Suffix:       file.FileExt,
			Size:         int(file.FileSize),
			Path:         filePath,
			Url:          uri,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err = m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[文件消息]"})
	}

	return err
}

// SendCode 代码消息
func (m *MessageService) SendCode(ctx context.Context, uid int, req *message.CodeMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeCode,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraCode{
			Lang: req.Lang,
			Code: req.Code,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[代码消息]"})
	}

	return err
}

// SendVote 投票消息
func (m *MessageService) SendVote(ctx context.Context, uid int, req *message.VoteMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   entity.ChatGroupMode,
		MsgType:    entity.MsgTypeVote,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	options := make(map[string]string)
	for i, value := range req.Options {
		options[fmt.Sprintf("%c", 65+i)] = value
	}

	num := m.groupMemberRepo.CountMemberTotal(ctx, int(req.Receiver.ReceiverId))

	err := m.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(data).Error; err != nil {
			return err
		}

		return tx.Create(&model.TalkRecordsVote{
			RecordId:     data.Id,
			UserId:       uid,
			Title:        req.Title,
			AnswerMode:   int(req.Mode),
			AnswerOption: jsonutil.Encode(options),
			AnswerNum:    int(num),
			IsAnonymous:  int(req.Anonymous),
		}).Error
	})

	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[投票消息]"})
	}

	return err
}

// SendEmoticon 表情消息
func (m *MessageService) SendEmoticon(ctx context.Context, uid int, req *message.EmoticonMessageRequest) error {

	emoticon := &model.EmoticonItem{}
	if err := m.Db().Model(&model.EmoticonItem{}).Where("id = ? and user_id = ?", req.EmoticonId, uid).First(emoticon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("表情信息不存在")
		}

		return err
	}

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeFile,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraFile{
			Type:         entity.GetMediaType(emoticon.FileSuffix),
			OriginalName: "图片表情",
			Suffix:       emoticon.FileSuffix,
			Size:         emoticon.FileSize,
			Path:         emoticon.Url,
			Url:          emoticon.Url,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[表情包]"})
	}

	return err
}

// SendForward 转发消息
func (m *MessageService) SendForward(ctx context.Context, uid int, req *message.ForwardMessageRequest) error {

	// 验证转发消息合法性
	if err := m.forward.Verify(ctx, uid, req); err != nil {
		return err
	}

	var items []*logic.ForwardRecord
	// 发送方式 1:逐条发送 2:合并发送
	if req.Mode == 1 {
		items, _ = m.forward.MultiSplitForward(ctx, uid, req)
	} else {
		items, _ = m.forward.MultiMergeForward(ctx, uid, req)
	}

	for _, item := range items {
		m.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(entity.MapStrAny{
			"event": entity.EventTalk,
			"data": jsonutil.Encode(entity.MapStrAny{
				"sender_id":   uid,
				"receiver_id": item.ReceiverId,
				"talk_type":   item.TalkType,
				"record_id":   item.RecordId,
			}),
		}))
	}

	return nil
}

// SendLocation 位置消息
func (m *MessageService) SendLocation(ctx context.Context, uid int, req *message.LocationMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.MsgTypeLocation,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraLocation{
			Longitude: req.Longitude,
			Latitude:  req.Latitude,
		}),
	}

	if req.Receiver.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, int(req.Receiver.ReceiverId))
	} else {
		data.Sequence = m.Sequence.Get(ctx, uid, int(req.Receiver.ReceiverId))
	}

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[位置消息]"})
	}

	return err
}

// SendBusinessCard 推送用户名片消息
func (m *MessageService) SendBusinessCard(ctx context.Context, uid int) error {
	panic("SendBusinessCard")
}

// SendLogin 推送用户登录消息
func (m *MessageService) SendLogin(ctx context.Context, uid int, req *message.LoginMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		Sequence:   m.Sequence.Get(ctx, 4257, uid),
		TalkType:   entity.ChatPrivateMode,
		MsgType:    entity.MsgTypeLogin,
		UserId:     4257, // 机器人ID
		ReceiverId: uid,
		Extra: jsonutil.Encode(&model.TalkRecordExtraLogin{
			IpAddress: req.Ip,
			Platform:  req.Platform,
			Agent:     req.Agent,
			Address:   req.Address,
			Reason:    req.Reason,
			Datetime:  timeutil.DateTime(),
		}),
	}

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[登录消息]"})
	}

	return err
}

// 发送消息后置处理
func (m *MessageService) afterHandle(ctx context.Context, record *model.TalkRecords, opts map[string]string) {

	if record.TalkType == entity.ChatPrivateMode {
		m.unreadStorage.Incr(ctx, entity.ChatPrivateMode, record.UserId, record.ReceiverId)

		if record.MsgType == entity.MsgTypeSystemText {
			m.unreadStorage.Incr(ctx, 1, record.ReceiverId, record.UserId)
		}
	} else if record.TalkType == entity.ChatGroupMode {

		// todo 需要加缓存
		ids := m.groupMemberRepo.GetMemberIds(ctx, record.ReceiverId)
		for _, uid := range ids {

			if uid == record.UserId {
				continue
			}

			m.unreadStorage.Incr(ctx, entity.ChatGroupMode, record.ReceiverId, uid)
		}
	}

	_ = m.messageStorage.Set(ctx, record.TalkType, record.UserId, record.ReceiverId, &cache.LastCacheMessage{
		Content:  opts["text"],
		Datetime: timeutil.DateTime(),
	})

	content := jsonutil.Encode(map[string]any{
		"event": entity.EventTalk,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"talk_type":   record.TalkType,
			"record_id":   record.Id,
		}),
	})

	// 点对点消息采用精确投递
	if record.TalkType == entity.ChatPrivateMode {
		sids := m.sidStorage.All(ctx, 1)

		// 小于三台服务器则采用全局广播
		if len(sids) <= 3 {
			if err := m.Redis().Publish(ctx, entity.ImTopicChat, content).Err(); err != nil {
				logger.Error(fmt.Sprintf("[ALL]消息推送失败 %s", err.Error()))
			}

			return
		}

		for _, sid := range m.sidStorage.All(ctx, 1) {
			for _, uid := range []int{record.UserId, record.ReceiverId} {
				if !m.clientStorage.IsCurrentServerOnline(ctx, sid, entity.ImChannelChat, strconv.Itoa(uid)) {
					continue
				}

				if err := m.Redis().Publish(ctx, fmt.Sprintf(entity.ImTopicChatPrivate, sid), content).Err(); err != nil {
					logger.WithFields(entity.H{
						"sid": sid,
					}).Error(fmt.Sprintf("[Private]消息推送失败 %s", err.Error()))
				}
			}
		}
	} else {
		if err := m.Redis().Publish(ctx, entity.ImTopicChat, content).Err(); err != nil {
			logger.Error(fmt.Sprintf("[ALL]消息推送失败 %s", err.Error()))
		}
	}
}
