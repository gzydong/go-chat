package service

import (
	"context"
	"errors"
	"fmt"
	"html"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
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

type IMessageService interface {
	// SendSystemText 系统文本消息
	SendSystemText(ctx context.Context, uid int, req *message.TextMessageRequest) error
	// SendText 文本消息
	SendText(ctx context.Context, uid int, req *message.TextMessageRequest) error
	// SendImage 图片文件消息
	SendImage(ctx context.Context, uid int, req *message.ImageMessageRequest) error
	// SendVoice 语音文件消息
	SendVoice(ctx context.Context, uid int, req *message.VoiceMessageRequest) error
	// SendVideo 视频文件消息
	SendVideo(ctx context.Context, uid int, req *message.VideoMessageRequest) error
	// SendFile 文件消息
	SendFile(ctx context.Context, uid int, req *message.FileMessageRequest) error
	// SendCode 代码消息
	SendCode(ctx context.Context, uid int, req *message.CodeMessageRequest) error
	// SendVote 投票消息
	SendVote(ctx context.Context, uid int, req *message.VoteMessageRequest) error
	// SendEmoticon 表情消息
	SendEmoticon(ctx context.Context, uid int, req *message.EmoticonMessageRequest) error
	// SendForward 转发消息
	SendForward(ctx context.Context, uid int, req *message.ForwardMessageRequest) error
	// SendLocation 位置消息
	SendLocation(ctx context.Context, uid int, req *message.LocationMessageRequest) error
	// SendBusinessCard 推送用户名片消息
	SendBusinessCard(ctx context.Context, uid int, req *message.CardMessageRequest) error
	// SendLogin 推送用户登录消息
	SendLogin(ctx context.Context, uid int, req *message.LoginMessageRequest) error
	// SendSysOther 推送其它消息
	SendSysOther(ctx context.Context, data *model.TalkRecords) error
	// SendMixedMessage 图文消息
	SendMixedMessage(ctx context.Context, uid int, req *message.MixedMessageRequest) error
}

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
	sequence            *repo.Sequence
	robotRepo           *repo.Robot
}

func NewMessageService(source *repo.Source, forward *logic.MessageForwardLogic, groupMemberRepo *repo.GroupMember, splitUploadRepo *repo.SplitUpload, talkRecordsVoteRepo *repo.TalkRecordsVote, fileSystem *filesystem.Filesystem, unreadStorage *cache.UnreadStorage, messageStorage *cache.MessageStorage, sidStorage *cache.ServerStorage, clientStorage *cache.ClientStorage, sequence *repo.Sequence, robotRepo *repo.Robot) *MessageService {
	return &MessageService{Source: source, forward: forward, groupMemberRepo: groupMemberRepo, splitUploadRepo: splitUploadRepo, talkRecordsVoteRepo: talkRecordsVoteRepo, fileSystem: fileSystem, unreadStorage: unreadStorage, messageStorage: messageStorage, sidStorage: sidStorage, clientStorage: clientStorage, sequence: sequence, robotRepo: robotRepo}
}

// SendSystemText 系统文本消息
func (m *MessageService) SendSystemText(ctx context.Context, uid int, req *message.TextMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgSysText,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Content:    html.EscapeString(req.Content),
	}

	return m.save(ctx, data)
}

// SendText 文本消息
func (m *MessageService) SendText(ctx context.Context, uid int, req *message.TextMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeText,
		QuoteId:    req.QuoteId,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Content:    strutil.EscapeHtml(req.Content),
	}

	return m.save(ctx, data)
}

// SendImage 图片文件消息
func (m *MessageService) SendImage(ctx context.Context, uid int, req *message.ImageMessageRequest) error {

	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeImage,
		QuoteId:    req.QuoteId,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraImage{
			Url:    req.Url,
			Width:  int(req.Width),
			Height: int(req.Height),
		}),
	}

	return m.save(ctx, data)
}

// SendVoice 语音文件消息
func (m *MessageService) SendVoice(ctx context.Context, uid int, req *message.VoiceMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeAudio,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraAudio{
			Suffix:   strutil.FileSuffix(req.Url),
			Size:     int(req.Size),
			Url:      req.Url,
			Duration: 0,
		}),
	}

	return m.save(ctx, data)
}

// SendVideo 视频文件消息
func (m *MessageService) SendVideo(ctx context.Context, uid int, req *message.VideoMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeVideo,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraVideo{
			Cover:    req.Cover,
			Size:     int(req.Size),
			Url:      req.Url,
			Duration: int(req.Duration),
		}),
	}

	return m.save(ctx, data)
}

// SendFile 文件消息
func (m *MessageService) SendFile(ctx context.Context, uid int, req *message.FileMessageRequest) error {

	file, err := m.splitUploadRepo.GetFile(ctx, uid, req.UploadId)
	if err != nil {
		return err
	}

	publicUrl := ""
	filePath := fmt.Sprintf("private/files/talks/%s/%s.%s", timeutil.DateNumber(), encrypt.Md5(strutil.Random(16)), file.FileExt)

	// 公开文件
	if entity.GetMediaType(file.FileExt) <= 3 {
		filePath = fmt.Sprintf("public/media/%s/%s.%s", timeutil.DateNumber(), encrypt.Md5(strutil.Random(16)), file.FileExt)
		publicUrl = m.fileSystem.Default.PublicUrl(filePath)
	}

	if err := m.fileSystem.Default.Copy(file.Path, filePath); err != nil {
		return err
	}

	data := &model.TalkRecords{
		MsgId:      encrypt.Md5(req.UploadId),
		TalkType:   int(req.Receiver.TalkType),
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
	}

	switch entity.GetMediaType(file.FileExt) {
	case entity.MediaFileAudio:
		data.MsgType = entity.ChatMsgTypeAudio
		data.Extra = jsonutil.Encode(&model.TalkRecordExtraAudio{
			Suffix:   file.FileExt,
			Size:     int(file.FileSize),
			Url:      publicUrl,
			Duration: 0,
		})
	case entity.MediaFileVideo:
		data.MsgType = entity.ChatMsgTypeVideo
		data.Extra = jsonutil.Encode(&model.TalkRecordExtraVideo{
			Cover:    "",
			Suffix:   file.FileExt,
			Size:     int(file.FileSize),
			Url:      publicUrl,
			Duration: 0,
		})
	case entity.MediaFileOther:
		data.MsgType = entity.ChatMsgTypeFile
		data.Extra = jsonutil.Encode(&model.TalkRecordExtraFile{
			Drive:  file.Drive,
			Name:   file.OriginalName,
			Suffix: file.FileExt,
			Size:   int(file.FileSize),
			Path:   filePath,
		})
	}

	return m.save(ctx, data)
}

// SendCode 代码消息
func (m *MessageService) SendCode(ctx context.Context, uid int, req *message.CodeMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeCode,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraCode{
			Lang: req.Lang,
			Code: req.Code,
		}),
	}

	return m.save(ctx, data)
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

	var emoticon model.EmoticonItem
	if err := m.Db().First(&emoticon, "id = ? and user_id = ?", req.EmoticonId, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("表情信息不存在")
		}

		return err
	}

	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeImage,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraImage{
			Url:    emoticon.Url,
			Width:  0,
			Height: 0,
		}),
	}

	return m.save(ctx, data)
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

	for _, record := range items {
		if record.TalkType == entity.ChatPrivateMode {
			m.unreadStorage.Incr(ctx, entity.ChatPrivateMode, uid, record.ReceiverId)
		} else if record.TalkType == entity.ChatGroupMode {
			pipe := m.Redis().Pipeline()
			for _, uid := range m.groupMemberRepo.GetMemberIds(ctx, record.ReceiverId) {
				m.unreadStorage.PipeIncr(ctx, pipe, entity.ChatGroupMode, record.ReceiverId, uid)
			}
			_, _ = pipe.Exec(ctx)
		}

		_ = m.messageStorage.Set(ctx, record.TalkType, uid, record.ReceiverId, &cache.LastCacheMessage{
			Content:  "[转发消息]",
			Datetime: timeutil.DateTime(),
		})
	}

	_, _ = m.Redis().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range items {
			data := jsonutil.Encode(map[string]any{
				"event": entity.SubEventImMessage,
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
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeLocation,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraLocation{
			Longitude:   req.Longitude,
			Latitude:    req.Latitude,
			Description: req.Description,
		}),
	}

	return m.save(ctx, data)
}

// SendBusinessCard 推送用户名片消息
func (m *MessageService) SendBusinessCard(ctx context.Context, uid int, req *message.CardMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeCard,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(&model.TalkRecordExtraCard{
			UserId: int(req.UserId),
		}),
	}

	return m.save(ctx, data)
}

// SendLogin 推送用户登录消息
func (m *MessageService) SendLogin(ctx context.Context, uid int, req *message.LoginMessageRequest) error {

	robot, err := m.robotRepo.GetLoginRobot(ctx)
	if err != nil {
		return err
	}

	data := &model.TalkRecords{
		TalkType:   entity.ChatPrivateMode,
		MsgType:    entity.ChatMsgTypeLogin,
		UserId:     robot.UserId,
		ReceiverId: uid,
		Extra: jsonutil.Encode(&model.TalkRecordExtraLogin{
			IP:       req.Ip,
			Platform: req.Platform,
			Agent:    req.Agent,
			Address:  req.Address,
			Reason:   req.Reason,
			Datetime: timeutil.DateTime(),
		}),
	}

	return m.save(ctx, data)
}

// SendMixedMessage 图文消息
func (m *MessageService) SendMixedMessage(ctx context.Context, uid int, req *message.MixedMessageRequest) error {

	items := make([]*model.TalkRecordExtraMixedItem, 0)

	for _, item := range req.Items {
		items = append(items, &model.TalkRecordExtraMixedItem{
			Type:    int(item.Type),
			Content: item.Content,
		})
	}

	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgTypeMixed,
		QuoteId:    req.QuoteId,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra:      jsonutil.Encode(model.TalkRecordExtraMixed{Items: items}),
	}

	return m.save(ctx, data)
}

// SendSysOther 推送其它消息
func (m *MessageService) SendSysOther(ctx context.Context, data *model.TalkRecords) error {
	return m.save(ctx, data)
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
		"event": entity.SubEventImMessageRevoke,
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

func (m *MessageService) save(ctx context.Context, data *model.TalkRecords) error {

	if data.MsgId == "" {
		data.MsgId = strutil.NewMsgId()
	}

	m.loadReply(ctx, data)

	m.loadSequence(ctx, data)

	if err := m.Db().WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	option := make(map[string]string)

	switch data.MsgType {
	case entity.ChatMsgTypeText:
		option["text"] = strutil.MtSubstr(strutil.ReplaceImgAll(data.Content), 0, 300)
	default:
		if value, ok := entity.ChatMsgTypeMapping[data.MsgType]; ok {
			option["text"] = value
		} else {
			option["text"] = "[未知消息]"
		}
	}

	m.afterHandle(ctx, data, option)

	return nil
}

func (m *MessageService) loadSequence(ctx context.Context, data *model.TalkRecords) {
	if data.TalkType == entity.ChatGroupMode {
		data.Sequence = m.sequence.Get(ctx, 0, data.ReceiverId)
	} else {
		data.Sequence = m.sequence.Get(ctx, data.UserId, data.ReceiverId)
	}
}

func (m *MessageService) loadReply(ctx context.Context, data *model.TalkRecords) {
	// 检测是否引用消息
	if data.QuoteId == "" {
		return
	}

	if data.Extra == "" {
		data.Extra = "{}"
	}

	extra := make(map[string]any)

	err := jsonutil.Decode(data.Extra, &extra)
	if err != nil {
		logger.Error("MessageService Json Decode err: ", err.Error())
		return
	}

	var record model.TalkRecords
	err = m.Db().Table("talk_records").Find(&record, "msg_id = ?", data.QuoteId).Error
	if err != nil {
		return
	}

	var user model.Users
	err = m.Db().Table("users").Select("nickname").Find(&user, "id = ?", record.UserId).Error
	if err != nil {
		return
	}

	reply := model.Reply{
		UserId:   record.UserId,
		Nickname: user.Nickname,
		MsgType:  1,
		Content:  record.Content,
		MsgId:    record.MsgId,
	}

	if record.MsgType != entity.ChatMsgTypeText {
		reply.Content = "[未知消息]"
		if value, ok := entity.ChatMsgTypeMapping[record.MsgType]; ok {
			reply.Content = value
		}
	}

	extra["reply"] = reply

	data.Extra = jsonutil.Encode(extra)
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
		"event": entity.SubEventImMessage,
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
