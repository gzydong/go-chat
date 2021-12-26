package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-chat/app/cache"
	"go-chat/app/dao"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/model"
	"go-chat/app/pkg/encrypt"
	"go-chat/app/pkg/filesystem"
	"go-chat/app/pkg/jsonutil"
	"go-chat/app/pkg/strutil"
	"go-chat/app/pkg/timeutil"
	"go-chat/app/pkg/utils"
	"go-chat/config"
	"gorm.io/gorm"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TalkMessageService struct {
	*BaseService
	config             *config.Config
	unreadTalkCache    *cache.UnreadTalkCache
	forwardService     *TalkMessageForwardService
	lastMessage        *cache.LastMessage
	talkRecordsVoteDao *dao.TalkRecordsVoteDao
	groupMemberDao     *dao.GroupMemberDao
	sidServer          *cache.SidServer
	client             *cache.WsClientSession
	fileSystem         *filesystem.Filesystem
}

func NewTalkMessageService(baseService *BaseService, config *config.Config, unreadTalkCache *cache.UnreadTalkCache, forwardService *TalkMessageForwardService, lastMessage *cache.LastMessage, talkRecordsVoteDao *dao.TalkRecordsVoteDao, groupMemberDao *dao.GroupMemberDao, sidServer *cache.SidServer, client *cache.WsClientSession, fileSystem *filesystem.Filesystem) *TalkMessageService {
	return &TalkMessageService{BaseService: baseService, config: config, unreadTalkCache: unreadTalkCache, forwardService: forwardService, lastMessage: lastMessage, talkRecordsVoteDao: talkRecordsVoteDao, groupMemberDao: groupMemberDao, sidServer: sidServer, client: client, fileSystem: fileSystem}
}

// SendTextMessage 发送文本消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendTextMessage(ctx context.Context, uid int, params *request.TextMessageRequest) error {
	record := &model.TalkRecords{
		TalkType:   params.TalkType,
		MsgType:    entity.MsgTypeText,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
		Content:    params.Text,
	}

	if err := s.db.Create(record).Error; err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": strutil.MtSubstr(&record.Content, 0, 30),
	})

	return nil
}

// SendCodeMessage 发送代码消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendCodeMessage(ctx context.Context, uid int, params *request.CodeMessageRequest) error {
	var (
		err    error
		record = &model.TalkRecords{
			TalkType:   params.TalkType,
			MsgType:    entity.MsgTypeCode,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}
	)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsCode{
			RecordId: record.Id,
			UserId:   uid,
			Lang:     params.Lang,
			Code:     params.Code,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": "[代码消息]",
	})

	return nil
}

// SendImageMessage 发送图片消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendImageMessage(ctx context.Context, uid int, params *request.ImageMessageRequest, file *multipart.FileHeader) error {
	var (
		err    error
		record = &model.TalkRecords{
			TalkType:   params.TalkType,
			MsgType:    entity.MsgTypeFile,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}
	)

	stream, err := filesystem.ReadMultipartStream(file)
	if err != nil {
		return err
	}

	ext := strutil.FileSuffix(file.Filename)
	m := utils.ReadFileImage(bytes.NewReader(stream))

	filePath := fmt.Sprintf("public/media/image/talk/%s/%s", timeutil.DateNumber(), strutil.GenImageName(ext, m.Width, m.Height))

	if err := s.fileSystem.Default.Write(stream, filePath); err != nil {
		logrus.Error("文件上传失败 err:", err.Error())
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsFile{
			RecordId:     record.Id,
			UserId:       uid,
			Source:       1,
			Type:         entity.GetMediaType(ext),
			Drive:        entity.FileDriveMode(s.fileSystem.Driver()),
			OriginalName: file.Filename,
			Suffix:       ext,
			Size:         int(file.Size),
			Path:         filePath,
			Url:          s.fileSystem.Default.PublicUrl(filePath),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": "[图片消息]",
	})

	return nil
}

// SendFileMessage 发送文件消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendFileMessage(ctx context.Context, uid int, params *request.FileMessageRequest, file *model.SplitUpload) error {

	var (
		err    error
		record = &model.TalkRecords{
			TalkType:   params.TalkType,
			MsgType:    entity.MsgTypeFile,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}
	)

	filePath := fmt.Sprintf("private/files/talks/%s/%s.%s", timeutil.DateNumber(), encrypt.Md5(strutil.Random(16)), file.FileExt)
	url := ""
	if entity.GetMediaType(file.FileExt) <= 3 {
		filePath = fmt.Sprintf("public/media/%s/%s.%s", timeutil.DateNumber(), encrypt.Md5(strutil.Random(16)), file.FileExt)
		url = s.fileSystem.Default.PublicUrl(filePath)
	}

	if err := s.fileSystem.Default.Copy(file.SaveDir, filePath); err != nil {
		logrus.Error("文件拷贝失败 err: ", err.Error())
		return err
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsFile{
			RecordId:     record.Id,
			UserId:       uid,
			Source:       1,
			Type:         entity.GetMediaType(file.FileExt),
			Drive:        file.Drive,
			OriginalName: file.OriginalName,
			Suffix:       file.FileExt,
			Size:         int(file.FileSize),
			Path:         filePath,
			Url:          url,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": "[文件消息]",
	})

	return nil
}

// SendCardMessage 发送用户名片消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendCardMessage(ctx context.Context, params *request.CardMessageRequest) {

}

// SendVoteMessage 发送投票消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendVoteMessage(ctx context.Context, uid int, params *request.VoteMessageRequest) error {
	var (
		err    error
		record = &model.TalkRecords{
			TalkType:   entity.GroupChat,
			MsgType:    entity.MsgTypeVote,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}
	)

	options := make(map[string]string)
	for i, value := range params.Options {
		options[fmt.Sprintf("%c", 65+i)] = value
	}

	num := s.groupMemberDao.CountMemberTotal(params.ReceiverId)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsVote{
			RecordId:     record.Id,
			UserId:       uid,
			Title:        params.Title,
			AnswerMode:   params.Mode,
			AnswerOption: jsonutil.JsonEncode(options),
			AnswerNum:    int(num),
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": "[投票消息]",
	})

	return nil
}

// SendEmoticonMessage 发送表情包消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendEmoticonMessage(ctx context.Context, uid int, params *request.EmoticonMessageRequest) error {
	var (
		err      error
		emoticon model.EmoticonItem
		record   = &model.TalkRecords{
			TalkType:   params.TalkType,
			MsgType:    entity.MsgTypeFile,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}
	)

	if err = s.db.Model(&model.EmoticonItem{}).Where("id = ?", params.EmoticonId).First(&emoticon).Error; err != nil {
		return err
	}

	if emoticon.UserId > 0 && emoticon.UserId != uid {
		return errors.New("表情包不存在！")
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsFile{
			RecordId:     record.Id,
			UserId:       uid,
			Source:       2,
			Type:         entity.GetMediaType(emoticon.FileSuffix),
			OriginalName: "图片表情",
			Suffix:       emoticon.FileSuffix,
			Size:         emoticon.FileSize,
			Path:         emoticon.Url,
			Url:          emoticon.Url,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": "[图片消息]",
	})

	return nil
}

// SendLocationMessage 发送位置消息
// @params uid     用户ID
// @params params  请求参数
func (s *TalkMessageService) SendLocationMessage(ctx context.Context, uid int, params *request.LocationMessageRequest) error {

	var (
		err    error
		record = &model.TalkRecords{
			TalkType:   params.TalkType,
			MsgType:    entity.MsgTypeLocation,
			UserId:     uid,
			ReceiverId: params.ReceiverId,
		}
	)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsLocation{
			RecordId:  record.Id,
			UserId:    uid,
			Longitude: params.Longitude,
			Latitude:  params.Latitude,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	s.afterHandle(ctx, record, map[string]string{
		"text": "[位置消息]",
	})

	return nil
}

// SendLocationMessage 撤销推送消息
// @params uid       用户ID
// @params recordId  消息记录ID
func (s *TalkMessageService) SendRevokeRecordMessage(ctx context.Context, uid int, recordId int) error {
	var (
		err    error
		record model.TalkRecords
	)

	if err = s.db.First(&record, recordId).Error; err != nil {
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

	if err = s.db.Model(&model.TalkRecords{Id: recordId}).Update("is_revoke", 1).Error; err != nil {
		return err
	}

	body := map[string]interface{}{
		"event": entity.EventRevokeTalk,
		"data": jsonutil.JsonEncode(map[string]interface{}{
			"record_id": record.Id,
		}),
	}

	s.rds.Publish(ctx, entity.IMGatewayAll, jsonutil.JsonEncode(body))

	return nil
}

// VoteHandle 投票处理
// @params uid       用户ID
// @params recordId  消息记录ID
func (s *TalkMessageService) VoteHandle(ctx context.Context, uid int, params *request.VoteMessageHandleRequest) (int, error) {
	var (
		err  error
		vote *model.QueryVoteModel
	)

	tx := s.db.Table("talk_records")
	tx.Select([]string{
		"talk_records.receiver_id", "talk_records.talk_type", "talk_records.msg_type",
		"vote.id as vote_id", "vote.id as record_id", "vote.answer_mode", "vote.answer_option",
		"vote.answer_num", "vote.status as vote_status",
	})
	tx.Joins("left join talk_records_vote as vote on vote.record_id = talk_records.id")
	tx.Where("talk_records.id = ?", params.RecordId)

	res := tx.Take(&vote)
	if err := res.Error; err != nil {
		return 0, err
	}

	if res.RowsAffected == 0 {
		return 0, fmt.Errorf("投票信息不存在[%d]", params.RecordId)
	}

	if vote.MsgType != entity.MsgTypeVote {
		return 0, fmt.Errorf("当前记录属于投票信息[%d]", vote.MsgType)
	}

	// 判断是否有投票权限

	var count int64
	s.db.Table("talk_records_vote_answer").Where("vote_id = ? and user_id = ？", vote.VoteId, uid).Count(&count)
	if count > 0 { // 判断是否已投票
		return 0, fmt.Errorf("不能重复投票[%d]", vote.VoteId)
	}

	options := strings.Split(params.Options, ",")
	sort.Strings(options)

	var answerOptions map[string]interface{}
	if err = jsonutil.JsonDecode(vote.AnswerOption, &answerOptions); err != nil {
		return 0, err
	}

	for _, option := range options {
		if _, ok := answerOptions[option]; !ok {
			return 0, fmt.Errorf("的投票选项不存在[%s]", option)
		}
	}

	// 判断是否单选
	if vote.AnswerMode == 0 {
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

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Table("talk_records_vote").Where("id = ?", vote.VoteId).Updates(map[string]interface{}{
			"answered_num": gorm.Expr("answered_num + 1"),
			"status":       gorm.Expr("if(answered_num >= answer_num, 1, 0)"),
		}).Error; err != nil {
			return err
		}

		if err = tx.Create(answers).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	_, _ = s.talkRecordsVoteDao.SetVoteAnswerUser(ctx, vote.VoteId)
	_, _ = s.talkRecordsVoteDao.SetVoteStatistics(ctx, vote.VoteId)

	return vote.VoteId, nil
}

type LoginInfo struct {
	UserId   int    `json:"user_id"`
	Ip       string `json:"ip"`
	Address  string `json:"address"`
	Platform string `json:"platform"`
	Agent    string `json:"agent"`
}

func (s *TalkMessageService) SendLoginMessage(ctx context.Context, login *LoginInfo) error {
	var (
		err    error
		record = &model.TalkRecords{
			TalkType:   entity.PrivateChat,
			MsgType:    entity.MsgTypeLogin,
			UserId:     4257,
			ReceiverId: login.UserId,
		}
	)

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err = s.db.Create(record).Error; err != nil {
			return err
		}

		if err = s.db.Create(&model.TalkRecordsLogin{
			RecordId: record.Id,
			UserId:   login.UserId,
			Ip:       login.Ip,
			Platform: login.Platform,
			Agent:    login.Agent,
			Address:  login.Address,
			Reason:   "常用设备登录",
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err == nil {
		s.afterHandle(ctx, record, map[string]string{
			"text": "[系统通知]",
		})
	}

	return err
}

// 发送消息后置处理
func (s *TalkMessageService) afterHandle(ctx context.Context, record *model.TalkRecords, opts map[string]string) {
	if record.TalkType == entity.PrivateChat {
		s.unreadTalkCache.Increment(ctx, record.UserId, record.ReceiverId)
	}

	_ = s.lastMessage.Set(ctx, record.TalkType, record.UserId, record.ReceiverId, &cache.LastCacheMessage{
		Content:  opts["text"],
		Datetime: timeutil.DateTime(),
	})

	// 推送消息至 redis

	body := map[string]interface{}{
		"event": entity.EventTalk,
		"data": jsonutil.JsonEncode(map[string]interface{}{
			"sender_id":   int64(record.UserId),
			"receiver_id": int64(record.ReceiverId),
			"talk_type":   record.TalkType,
			"record_id":   int64(record.Id),
		}),
	}

	content := jsonutil.JsonEncode(body)

	// 点对点消息采用精确投递
	if record.TalkType == entity.PrivateChat {
		sids := s.sidServer.GetServerAll(ctx, 1)

		// 小于两台服务器则采用全局广播
		if len(sids) <= 3 {
			s.rds.Publish(ctx, entity.IMGatewayAll, content)
		} else {
			to := []int{record.UserId, record.ReceiverId}
			for _, sid := range s.sidServer.GetServerAll(ctx, 1) {
				for _, uid := range to {
					if s.client.IsCurrentServerOnline(ctx, sid, entity.ImChannelDefault, strconv.Itoa(uid)) {
						s.rds.Publish(ctx, entity.GetIMGatewayPrivate(sid), content)
					}
				}
			}
		}
	} else {
		s.rds.Publish(ctx, entity.IMGatewayAll, content)
	}
}
