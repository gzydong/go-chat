package service

import (
	"go-chat/app/http/request"
	"go-chat/app/repository"
	"go-chat/config"
)

type TalkMessageService struct {
	config               *config.Config
	talkRecordsRepo      *repository.TalkRecordsRepo
	talkRecordsCodeRepo  *repository.TalkRecordsCodeRepo
	talkRecordsLoginRepo *repository.TalkRecordsLoginRepo
	talkRecordsFileRepo  *repository.TalkRecordsFileRepo
	talkRecordsVoteRepo  *repository.TalkRecordsVoteRepo
}

func NewTalkMessageService(
	config *config.Config,
	talkRecordsRepo *repository.TalkRecordsRepo,
	talkRecordsCodeRepo *repository.TalkRecordsCodeRepo,
	talkRecordsLoginRepo *repository.TalkRecordsLoginRepo,
	talkRecordsFileRepo *repository.TalkRecordsFileRepo,
	talkRecordsVoteRepo *repository.TalkRecordsVoteRepo,
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

func (s *TalkMessageService) SendTextMessage(params *request.TextMessageRequest) {

}

func (s *TalkMessageService) SendCodeMessage(params *request.CodeMessageRequest) {

}

func (s *TalkMessageService) SendImageMessage(params *request.ImageMessageRequest) {

}

func (s *TalkMessageService) SendFileMessage(params *request.FileMessageRequest) {

}

func (s *TalkMessageService) SendCardMessage(params *request.CardMessageRequest) {

}

func (s *TalkMessageService) SendVoteMessage(params *request.VoteMessageRequest) {

}

func (s *TalkMessageService) SendEmoticonMessage(params *request.EmoticonMessageRequest) {

}

func (s *TalkMessageService) SendForwardMessage(params *request.ForwardMessageRequest) {

}
