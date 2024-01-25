package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

var _ IMessageService = (*MessageService)(nil)

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
	// Revoke 撤回消息
	Revoke(ctx context.Context, uid int, msgId string) error
	// Vote 投票
	Vote(ctx context.Context, uid int, msgId string, optionsValue string) (*repo.VoteStatistics, error)
}

type MessageService struct {
	*repo.Source
	MessageForwardLogic *logic.MessageForwardLogic
	GroupMemberRepo     *repo.GroupMember
	SplitUploadRepo     *repo.SplitUpload
	TalkRecordsVoteRepo *repo.TalkRecordsVote
	Filesystem          filesystem.IFilesystem
	UnreadStorage       *cache.UnreadStorage
	MessageStorage      *cache.MessageStorage
	ServerStorage       *cache.ServerStorage
	ClientStorage       *cache.ClientStorage
	Sequence            *repo.Sequence
	RobotRepo           *repo.Robot
}

// SendSystemText 系统文本消息
func (m *MessageService) SendSystemText(ctx context.Context, uid int, req *message.TextMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   int(req.Receiver.TalkType),
		MsgType:    entity.ChatMsgSysText,
		UserId:     uid,
		ReceiverId: int(req.Receiver.ReceiverId),
		Extra: jsonutil.Encode(model.TalkRecordExtraText{
			Content: strutil.EscapeHtml(req.Content),
		}),
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
		Extra: jsonutil.Encode(model.TalkRecordExtraText{
			Content:  strutil.EscapeHtml(req.Content),
			Mentions: req.Mentions,
		}),
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
	now := time.Now()

	file, err := m.SplitUploadRepo.GetFile(ctx, uid, req.UploadId)
	if err != nil {
		return err
	}

	publicUrl := ""
	filePath := fmt.Sprintf("talk-files/%s/%s.%s", now.Format("200601"), uuid.New().String(), file.FileExt)

	// 公开文件
	if entity.GetMediaType(file.FileExt) <= 3 {
		filePath = strutil.GenMediaObjectName(file.FileExt, 0, 0)
		// 如果是多媒体文件，则将私有文件转移到公开文件
		if err := m.Filesystem.CopyObject(
			m.Filesystem.BucketPrivateName(), file.Path,
			m.Filesystem.BucketPublicName(), filePath,
		); err != nil {
			return err
		}

		publicUrl = m.Filesystem.PublicUrl(m.Filesystem.BucketPublicName(), filePath)
	} else {
		if err := m.Filesystem.Copy(m.Filesystem.BucketPrivateName(), file.Path, filePath); err != nil {
			return err
		}
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
			Size:     int(file.FileSize),
			Url:      publicUrl,
			Duration: 0,
		})
	case entity.MediaFileVideo:
		data.MsgType = entity.ChatMsgTypeVideo
		data.Extra = jsonutil.Encode(&model.TalkRecordExtraVideo{
			Cover:    "",
			Size:     int(file.FileSize),
			Url:      publicUrl,
			Duration: 0,
		})
	case entity.MediaFileOther:
		data.MsgType = entity.ChatMsgTypeFile
		data.Extra = jsonutil.Encode(&model.TalkRecordExtraFile{
			Drive: file.Drive,
			Name:  file.OriginalName,
			Size:  int(file.FileSize),
			Path:  filePath,
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

	num := m.GroupMemberRepo.CountMemberTotal(ctx, int(req.Receiver.ReceiverId))

	err := m.Source.Db().WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		if err := tx.Create(data).Error; err != nil {
			return err
		}

		return tx.Create(&model.TalkRecordsVote{
			MsgId:        data.MsgId,
			UserId:       uid,
			Title:        req.Title,
			AnswerMode:   int(req.Mode),
			AnswerOption: jsonutil.Encode(options),
			AnswerNum:    int(num),
			IsAnonymous:  int(req.Anonymous),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}).Error
	})

	if err == nil {
		m.afterHandle(ctx, data, entity.TalkLastMessage{
			MsgId:      data.MsgId,
			Sequence:   int(data.Sequence),
			MsgType:    data.MsgType,
			UserId:     data.UserId,
			ReceiverId: data.ReceiverId,
			CreatedAt:  time.Now().Format(time.DateTime),
			Content:    "投票消息",
		})
	}

	return err
}

// SendEmoticon 表情消息
func (m *MessageService) SendEmoticon(ctx context.Context, uid int, req *message.EmoticonMessageRequest) error {

	var emoticon model.EmoticonItem
	if err := m.Source.Db().First(&emoticon, "id = ? and user_id = ?", req.EmoticonId, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
	if err := m.MessageForwardLogic.Verify(ctx, uid, req); err != nil {
		return err
	}

	var (
		err   error
		items []*logic.ForwardRecord
	)

	// 发送方式 1:逐条发送 2:合并发送
	if req.Mode == 1 {
		items, err = m.MessageForwardLogic.MultiSplitForward(ctx, uid, req)
	} else {
		items, err = m.MessageForwardLogic.MultiMergeForward(ctx, uid, req)
	}

	if err != nil {
		return err
	}

	for _, record := range items {
		if record.TalkType == entity.ChatPrivateMode {
			m.UnreadStorage.Incr(ctx, entity.ChatPrivateMode, uid, record.ReceiverId)
		} else if record.TalkType == entity.ChatGroupMode {
			pipe := m.Source.Redis().Pipeline()
			for _, uid := range m.GroupMemberRepo.GetMemberIds(ctx, record.ReceiverId) {
				m.UnreadStorage.PipeIncr(ctx, pipe, entity.ChatGroupMode, record.ReceiverId, uid)
			}
			_, _ = pipe.Exec(ctx)
		}

		_ = m.MessageStorage.Set(ctx, record.TalkType, uid, record.ReceiverId, &cache.LastCacheMessage{
			Content:  "[转发消息]",
			Datetime: timeutil.DateTime(),
		})
	}

	_, _ = m.Source.Redis().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, item := range items {
			data := jsonutil.Encode(map[string]any{
				"event": entity.SubEventImMessage,
				"data": jsonutil.Encode(map[string]any{
					"sender_id":   uid,
					"receiver_id": item.ReceiverId,
					"talk_type":   item.TalkType,
					"msg_id":      item.MsgId,
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

	robot, err := m.RobotRepo.GetLoginRobot(ctx)
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
func (m *MessageService) Revoke(ctx context.Context, uid int, msgId string) error {
	var record model.TalkRecords
	if err := m.Source.Db().First(&record, "msg_id = ?", msgId).Error; err != nil {
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

	if err := m.Source.Db().Model(&model.TalkRecords{Id: record.Id}).Update("is_revoke", 1).Error; err != nil {
		return err
	}

	var user model.Users
	if err := m.Db().WithContext(ctx).Select("id,nickname").First(&user, record.UserId).Error; err != nil {
		return err
	}

	_ = m.MessageStorage.Set(ctx, record.TalkType, record.UserId, record.ReceiverId, &cache.LastCacheMessage{
		Content:  fmt.Sprintf("%s: 撤回了一条消息", user.Nickname),
		Datetime: timeutil.DateTime(),
	})

	body := map[string]any{
		"event": entity.SubEventImMessageRevoke,
		"data": jsonutil.Encode(map[string]any{
			"msg_id": record.MsgId,
		}),
	}

	m.Source.Redis().Publish(ctx, entity.ImTopicChat, jsonutil.Encode(body))

	return nil
}

// Vote 投票
func (m *MessageService) Vote(ctx context.Context, uid int, msgId string, optionsValue string) (*repo.VoteStatistics, error) {

	db := m.Source.Db().WithContext(ctx)

	query := db.Table("talk_records")
	query.Select([]string{
		"talk_records.receiver_id", "talk_records.talk_type", "talk_records.msg_type",
		"vote.id as vote_id", "vote.msg_id as msg_id", "vote.answer_mode", "vote.answer_option",
		"vote.answer_num", "vote.status as vote_status",
	})
	query.Joins("left join talk_records_vote as vote on vote.msg_id = talk_records.msg_id")
	query.Where("talk_records.msg_id = ?", msgId)

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

	err := m.Source.Db().Transaction(func(tx *gorm.DB) error {
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

	_, _ = m.TalkRecordsVoteRepo.SetVoteAnswerUser(ctx, vote.VoteId)
	_, _ = m.TalkRecordsVoteRepo.SetVoteStatistics(ctx, vote.VoteId)
	info, _ := m.TalkRecordsVoteRepo.GetVoteStatistics(ctx, vote.VoteId)

	return info, nil
}

func (m *MessageService) save(ctx context.Context, data *model.TalkRecords) error {

	if data.MsgId == "" {
		data.MsgId = strutil.NewMsgId()
	}

	m.loadReply(ctx, data)

	m.loadSequence(ctx, data)

	if err := m.Source.Db().WithContext(ctx).Create(data).Error; err != nil {
		return err
	}

	lastMessage := entity.TalkLastMessage{
		MsgId:      data.MsgId,
		Sequence:   int(data.Sequence),
		MsgType:    data.MsgType,
		UserId:     data.UserId,
		ReceiverId: data.ReceiverId,
		CreatedAt:  time.Now().Format(time.DateTime),
	}

	switch data.MsgType {
	case entity.ChatMsgTypeText:
		extra := model.TalkRecordExtraText{}
		if err := jsonutil.Decode(data.Extra, &extra); err != nil {
			logger.Errorf("MessageService Json Decode err: %s", err.Error())
			return err
		}

		lastMessage.Content = strutil.MtSubstr(strutil.ReplaceImgAll(extra.Content), 0, 300)
	default:
		if value, ok := entity.ChatMsgTypeMapping[data.MsgType]; ok {
			lastMessage.Content = value
		} else {
			lastMessage.Content = "[未知消息]"
		}
	}

	m.afterHandle(ctx, data, lastMessage)

	return nil
}

func (m *MessageService) loadSequence(ctx context.Context, data *model.TalkRecords) {
	if data.TalkType == entity.ChatGroupMode {
		data.Sequence = m.Sequence.Get(ctx, 0, data.ReceiverId)
	} else {
		data.Sequence = m.Sequence.Get(ctx, data.UserId, data.ReceiverId)
	}
}

func (m *MessageService) loadReply(_ context.Context, data *model.TalkRecords) {
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
		logger.Errorf("MessageService Json Decode err: %s", err.Error())
		return
	}

	var record model.TalkRecords
	err = m.Source.Db().Table("talk_records").Find(&record, "msg_id = ?", data.QuoteId).Error
	if err != nil {
		return
	}

	var user model.Users
	err = m.Source.Db().Table("users").Select("nickname").Find(&user, "id = ?", record.UserId).Error
	if err != nil {
		return
	}

	reply := model.Reply{
		UserId:   record.UserId,
		Nickname: user.Nickname,
		MsgType:  1,
		MsgId:    record.MsgId,
	}

	if record.MsgType != entity.ChatMsgTypeText {
		reply.Content = "[未知消息]"
		if value, ok := entity.ChatMsgTypeMapping[record.MsgType]; ok {
			reply.Content = value
		}
	} else {
		extra := model.TalkRecordExtraText{}
		if err := jsonutil.Decode(record.Extra, &extra); err != nil {
			logger.Errorf("loadReply Json Decode err: %s", err.Error())
			return
		}

		reply.Content = extra.Content
	}

	extra["reply"] = reply

	data.Extra = jsonutil.Encode(extra)
}

// 发送消息后置处理
func (m *MessageService) afterHandle(ctx context.Context, record *model.TalkRecords, opt entity.TalkLastMessage) {

	if record.TalkType == entity.ChatPrivateMode {
		m.UnreadStorage.Incr(ctx, entity.ChatPrivateMode, record.UserId, record.ReceiverId)
		if record.MsgType == entity.ChatMsgSysText {
			m.UnreadStorage.Incr(ctx, 1, record.ReceiverId, record.UserId)
		}
	} else if record.TalkType == entity.ChatGroupMode {
		pipe := m.Source.Redis().Pipeline()
		for _, uid := range m.GroupMemberRepo.GetMemberIds(ctx, record.ReceiverId) {
			if uid != record.UserId {
				m.UnreadStorage.PipeIncr(ctx, pipe, entity.ChatGroupMode, record.ReceiverId, uid)
			}
		}
		_, _ = pipe.Exec(ctx)
	}

	_ = m.MessageStorage.Set(ctx, record.TalkType, record.UserId, record.ReceiverId, &cache.LastCacheMessage{
		Content:  opt.Content,
		Datetime: opt.CreatedAt,
	})

	content := jsonutil.Encode(map[string]any{
		"event": entity.SubEventImMessage,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserId,
			"receiver_id": record.ReceiverId,
			"talk_type":   record.TalkType,
			"msg_id":      record.MsgId,
		}),
	})

	if record.TalkType == entity.ChatPrivateMode {
		sids := m.ServerStorage.All(ctx, 1)

		if len(sids) > 3 {
			pipe := m.Source.Redis().Pipeline()

			for _, sid := range sids {
				for _, uid := range []int{record.UserId, record.ReceiverId} {
					if !m.ClientStorage.IsCurrentServerOnline(ctx, sid, entity.ImChannelChat, strconv.Itoa(uid)) {
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

	if err := m.Source.Redis().Publish(ctx, entity.ImTopicChat, content).Err(); err != nil {
		logger.Errorf("[ALL]消息推送失败 %s", err.Error())
	}
}
