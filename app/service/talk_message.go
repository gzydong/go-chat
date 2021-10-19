package service

import (
	"context"
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
