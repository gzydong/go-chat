package v1

import "go-chat/app/service"

type Talk struct {
	talkService *service.TalkService
}

func NewTalkHandler(
	talk *service.TalkService,
) *Talk {
	return &Talk{
		talkService: talk,
	}
}
