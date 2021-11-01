package service

import (
	"context"
)

type TalkMessageForwardService struct {
	*BaseService
}

func NewTalkMessageForwardService(base *BaseService) *TalkMessageForwardService {
	return &TalkMessageForwardService{base}
}

// 验证消息转发
func (t *TalkMessageForwardService) verifyForward() {

}

// MultiMergeForward 转发消息（多条合并转发）
func (t TalkMessageForwardService) MultiMergeForward(ctx context.Context, uid int, receiverId int, talkType int, recordsIds []int, receives []int) {

}

//

// MultiSplitForward 转发消息（多条拆分转发）
func (t TalkMessageForwardService) MultiSplitForward(ctx context.Context, uid int, receiverId int, talkType int, recordsIds []int, receives []int) {
}
