package service

import (
	"context"
	"go-chat/app/dao"
	"go-chat/app/http/request"
	"go-chat/config"
)

type TalkMessageService struct {
	config               *config.Config
	talkRecordsRepo      *dao.TalkRecordsDao
	talkRecordsCodeRepo  *dao.TalkRecordsCodeDao
	talkRecordsLoginRepo *dao.TalkRecordsLoginDao
	talkRecordsFileRepo  *dao.TalkRecordsFileDao
	talkRecordsVoteRepo  *dao.TalkRecordsVoteDao
}

func NewTalkMessageService(
	config *config.Config,
	talkRecordsRepo *dao.TalkRecordsDao,
	talkRecordsCodeRepo *dao.TalkRecordsCodeDao,
	talkRecordsLoginRepo *dao.TalkRecordsLoginDao,
	talkRecordsFileRepo *dao.TalkRecordsFileDao,
	talkRecordsVoteRepo *dao.TalkRecordsVoteDao,
) *TalkMessageService {
	return &TalkMessageService{
		config:               config,
		talkRecordsRepo:      talkRecordsRepo,
		talkRecordsCodeRepo:  talkRecordsCodeRepo,
		talkRecordsLoginRepo: talkRecordsLoginRepo,
		talkRecordsFileRepo:  talkRecordsFileRepo,
		talkRecordsVoteRepo:  talkRecordsVoteRepo,
	}
}

func (s *TalkMessageService) SendTextMessage(ctx context.Context, params *request.TextMessageRequest) {

}

func (s *TalkMessageService) SendCodeMessage(ctx context.Context, params *request.CodeMessageRequest) {

}

func (s *TalkMessageService) SendImageMessage(ctx context.Context, params *request.ImageMessageRequest) {

}

func (s *TalkMessageService) SendFileMessage(ctx context.Context, params *request.FileMessageRequest) {

}

func (s *TalkMessageService) SendCardMessage(ctx context.Context, params *request.CardMessageRequest) {

}

func (s *TalkMessageService) SendVoteMessage(ctx context.Context, params *request.VoteMessageRequest) {

}

func (s *TalkMessageService) SendEmoticonMessage(ctx context.Context, params *request.EmoticonMessageRequest) {

}

func (s *TalkMessageService) SendForwardMessage(ctx context.Context, params *request.ForwardMessageRequest) {

}
