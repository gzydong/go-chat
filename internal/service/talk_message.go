package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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
	forward             *logic.MessageForwardLogic
	groupMemberRepo     *repo.GroupMember
	splitUploadRepo     *repo.SplitUpload
	talkRecordsVoteRepo *repo.TalkRecordsVote
	fileSystem          *filesystem.Filesystem
	unreadStorage       *cache.UnreadStorage
	messageStorage      *cache.MessageStorage
	sidStorage          *cache.ServerStorage
	clientStorage       *cache.ClientStorage
	Sequence            *repo.Sequence
}

func NewMessageService(source *repo.Source, forward *logic.MessageForwardLogic, groupMemberRepo *repo.GroupMember, splitUploadRepo *repo.SplitUpload, talkRecordsVoteRepo *repo.TalkRecordsVote, fileSystem *filesystem.Filesystem, unreadStorage *cache.UnreadStorage, messageStorage *cache.MessageStorage, sidStorage *cache.ServerStorage, clientStorage *cache.ClientStorage, sequence *repo.Sequence) *MessageService {
	return &MessageService{Source: source, forward: forward, groupMemberRepo: groupMemberRepo, splitUploadRepo: splitUploadRepo, talkRecordsVoteRepo: talkRecordsVoteRepo, fileSystem: fileSystem, unreadStorage: unreadStorage, messageStorage: messageStorage, sidStorage: sidStorage, clientStorage: clientStorage, Sequence: sequence}
}

// SendSystemText 系统文本消息
func (m *MessageService) SendSystemText(ctx context.Context, uid int, req *message.TextMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgSysText,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Content:    html.EscapeString(req.Content),
	}

	m.loadSequence(ctx, data)

	if err := m.Db().WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	m.afterHandle(ctx, data, map[string]string{
		"text": strutil.MtSubstr(data.Content, 0, 300),
	})

	return nil
}

// SendText 文本消息
func (m *MessageService) SendText(ctx context.Context, uid int, req *message.TextMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeText,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Content:    html.EscapeString(req.Content),
	}

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeFile,
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

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeFile,
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

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeFile,
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

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeFile,
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

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeCode,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraCode{
			Lang: req.Lang,
			Code: req.Code,
		}),
	}

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeVote,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
	}

	m.loadSequence(ctx, data)

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
		MsgType:    entity.ChatMsgTypeFile,
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

	m.loadSequence(ctx, data)

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

	var (
		err   error
		items []*logic.ForwardRecord
	)

	// 发送方式 1:逐条发送 2:合并发送
	if req.Mode == 1 {
		items, err = m.forward.MultiSplitForward(ctx, uid, req)
	} else {
		items, err = m.forward.MultiMergeForward(ctx, uid, req)
	}

	if err != nil {
		return err
	}

	_, _ = m.Redis().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range items {
			data := jsonutil.Encode(map[string]any{
				"event": entity.EventTalk,
				"data": jsonutil.Encode(map[string]any{
					"sender_id":   uid,
					"receiver_id": item.ReceiverId,
					"talk_type":   item.TalkType,
					"record_id":   item.RecordId,
				}),
			})

			pipe.Publish(ctx, entity.ImTopicChat, data)
		}
		return nil
	})

	return nil
}

// SendLocation 位置消息
func (m *MessageService) SendLocation(ctx context.Context, uid int, req *message.LocationMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeLocation,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraLocation{
			Longitude: req.Longitude,
			Latitude:  req.Latitude,
		}),
	}

	m.loadSequence(ctx, data)

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[位置消息]"})
	}

	return err
}

// SendBusinessCard 推送用户名片消息
func (m *MessageService) SendBusinessCard(ctx context.Context, uid int, req *message.CardMessageRequest) error {
	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeCard,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraCard{
			UserId: int(req.UserId),
		}),
	}

	m.loadSequence(ctx, data)

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[分享名片]"})
	}

	return err
}

// SendLogin 推送用户登录消息
func (m *MessageService) SendLogin(ctx context.Context, uid int, req *message.LoginMessageRequest) error {

	data := &model.TalkRecords{
		MsgId:      strutil.NewMsgId(),
		TalkType:   entity.ChatPrivateMode,
		MsgType:    entity.ChatMsgTypeLogin,
		UserId:     4257,
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

	m.loadSequence(ctx, data)

	err := m.Db().WithContext(ctx).Create(data).Error
	if err == nil {
		m.afterHandle(ctx, data, map[string]string{"text": "[登录消息]"})
	}

	return err
}

// Revoke 撤回消息
func (m *MessageService) Revoke(ctx context.Context, uid int, recordId int) error {

	var record model.TalkRecords
	if err := m.Db().First(&record, recordId).Error; err != nil {
		return err
	}

	if record.IsRevoke == 1 {
		return nil
	}

	if record.UserId != uid {
		return errors.New("无权撤回回消息")
	}

	if time.Now().Unix() > record.CreatedAt.Add(3*time.Minute).Unix() {
		return errors.New("超出有效撤回时间范围，无法进行撤销！")
	}

	if err := m.Db().Model(&model.TalkRecords{Id: recordId}).Update("is_revoke", 1).Error; err != nil {
		return err
	}

	body := map[string]any{
		"event": entity.EventTalkRevoke,
		"data": jsonutil.Encode(map[string]any{
			"record_id": record.Id,
		}),
	}

	m.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))

	return nil
}

// Vote 投票
func (m *MessageService) Vote(ctx context.Context, uid int, recordId int, optionsValue string) (*repo.VoteStatistics, error) {

	db := m.Db().WithContext(ctx)

	query := db.Table("talk_records")
	query.Select([]string{
		"talk_records.receiver_id", "talk_records.talk_type", "talk_records.msg_type",
		"vote.id as vote_id", "vote.id as record_id", "vote.answer_mode", "vote.answer_option",
		"vote.answer_num", "vote.status as vote_status",
	})
	query.Joins("left join talk_records_vote as vote on vote.record_id = talk_records.id")
	query.Where("talk_records.id = ?", recordId)

	var vote model.QueryVoteModel
	if err := query.Take(&vote).Error; err != nil {
		return nil, err
	}

	if vote.MsgType != entity.ChatMsgTypeVote {
		return nil, fmt.Errorf("当前记录不属于投票信息[%d]", vote.MsgType)
	}

	if vote.TalkType == entity.ChatGroupMode {
		var count int64
		db.Table("group_member").Where("group_id = ? and user_id = ? and is_quit = 0", vote.ReceiverId, uid).Count(&count)
		if count == 0 {
			return nil, errors.New("暂无投票权限！")
		}
	}

	var count int64
	db.Table("talk_records_vote_answer").Where("vote_id = ? and user_id = ？", vote.VoteId, uid).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("重复投票[%d]", vote.VoteId)
	}

	options := strings.Split(optionsValue, ",")
	sort.Strings(options)

	var answerOptions map[string]any
	if err := jsonutil.Decode(vote.AnswerOption, &answerOptions); err != nil {
		return nil, err
	}

	for _, option := range options {
		if _, ok := answerOptions[option]; !ok {
			return nil, fmt.Errorf("投票选项不合法[%s]", option)
		}
	}

	if vote.AnswerMode == model.VoteAnswerModeSingleChoice {
		options = options[:1]
	}

	answers := make([]*model.TalkRecordsVoteAnswer, 0, len(options))
	for _, option := range options {
		answers = append(answers, &model.TalkRecordsVoteAnswer{
			VoteId: vote.VoteId,
			UserId: uid,
			Option: option,
		})
	}

	err := m.Db().Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("talk_records_vote").Where("id = ?", vote.VoteId).Updates(map[string]any{
			"answered_num": gorm.Expr("answered_num + 1"),
			"status":       gorm.Expr("if(answered_num >= answer_num, 1, 0)"),
		}).Error; err != nil {
			return err
		}

		return tx.Create(answers).Error
	})

	if err != nil {
		return nil, err
	}

	_, _ = m.talkRecordsVoteRepo.SetVoteAnswerUser(ctx, vote.VoteId)
	_, _ = m.talkRecordsVoteRepo.SetVoteStatistics(ctx, vote.VoteId)
	info, _ := m.talkRecordsVoteRepo.GetVoteStatistics(ctx, vote.VoteId)

	return info, nil
}

// 发送消息后置处理
func (m *MessageService) afterHandle(ctx context.Context, record *model.TalkRecords, opt map[string]string) {

	if record.TalkType == entity.ChatPrivateMode {
		m.unreadStorage.Incr(ctx, entity.ChatPrivateMode, record.UserId, record.ReceiverId)
		if record.MsgType == entity.ChatMsgSysText {
			m.unreadStorage.Incr(ctx, 1, record.ReceiverId, record.UserId)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		pipe := m.Redis().Pipeline()
		for _, uid := range m.groupMemberRepo.GetMemberIds(ctx, record.ReceiverId) {
			if uid != record.UserId {
				m.unreadStorage.PipeIncr(ctx, pipe, entity.ChatGroupMode, record.ReceiverId, uid)
			}
		}
		_, _ = pipe.Exec(ctx)
	}

	_ = m.messageStorage.Set(ctx, record.TalkType, record.UserId, record.ReceiverId, &cache.LastCacheMessage{
		Content:  opt["text"],
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

	if record.TalkType == entity.ChatPrivateMode {
		sids := m.sidStorage.All(ctx, 1)

		if len(sids) > 3 {
			pipe := m.Redis().Pipeline()

			for _, sid := range sids {
				for _, uid := range []int{record.UserId, record.ReceiverId} {
					if !m.clientStorage.IsCurrentServerOnline(ctx, sid, entity.ImChannelChat, strconv.Itoa(uid)) {
						continue
					}

					pipe.Publish(ctx, fmt.Sprintf(entity.ImTopicChatPrivate, sid), content)
				}
			}

			if _, err := pipe.Exec(ctx); err == nil {
				return
			}
		}
	}

	if err := m.Redis().Publish(ctx, entity.ImTopicChat, content).Err(); err != nil {
		logger.Error(fmt.Sprintf("[ALL]消息推送失败 %s", err.Error()))
	}
}

func (m *MessageService) loadSequence(ctx context.Context, data *model.TalkRecords) {
	if data.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, data.ReceiverId)
	} else {
		data.Sequence = m.Sequence.Get(ctx, data.UserId, data.ReceiverId)
	}
}
